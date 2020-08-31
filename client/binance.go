package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
	"github.com/vietnamz/cli-common/api/types"
	"github.com/vietnamz/cli-common/daemon/openssl"
	"net/url"
	"os"
	"time"
)


var API_KEY = os.Getenv("BINANCE_KEY")
var API_SECRET = os.Getenv("BINANCE_SECRET")
const (
	PublicEndPoint =" test"
)
// BinanceEndPoint represents the general information for handling binance request.
// These step includes:
// + Take the binance server timestamp.
// + convert the request body or query string into a signature.
type BinanceEndPoint struct {
	Key string
	signature string
	ServerTime int64
}


// Binance network has the format [verb]Binance[xx]
// where:
//		verb: get, create, update, delete, ping.
//		xx  : resource.
func (cli *Client) PingBinanceServer( ctx context.Context) bool {
	resp, err := cli.get(ctx, "/api/v3/ping", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		logrus.Errorf("failed to ping to binance server %s", err)
		return false
	}
	return true
}

func (cli *Client) GetBinanceServerTime( ctx context.Context) (int64, error) {
	var s types.BinanceServerTime
	resp, err := cli.get(ctx, "/api/v3/time", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return 0, errors.New("failed to get server time")
	}
	logrus.Errorf("status %d", resp.statusCode)
	err = json.NewDecoder(resp.body).Decode(&s)
	if err != nil {
		logrus.Errorf("failed to decode json with err %s", err)
	}
	return s.ServerTime, err
}


/*  query an order.
	params:
		ctx: the context,
		symbol: Mandatory, return an error if not exist.
		orderId->recvWindow: optional.
		Optional param should be filled with default value.
		Ex: orderId is int64. So we can fill zero.
			origClientOrderId is string. So we can fill "".
			it will be omitted as define in the type.
		return an order response along with the error.
*/
func (cli *Client) GetOrder(ctx context.Context,
							symbol string, orderId int64,
							origClientOrderId string,
							recvWindow int64) (*types.BinanceOrder,error){
	var orders types.BinanceOrder
	params := types.BinanceOrderParams {
		Symbol: symbol,
		OrderId: orderId,
		OrigClientOrderId: origClientOrderId,
		RecvWindow: recvWindow,
	}
	v, _ := query.Values(params)
	// take the unit time in milliseconds.
	params.Timestamp = time.Now().Unix() * 1000

	// encode the golang type to query string.
	// Note: excluding the ? symbol.
	queryParams := v.Encode()

	// sign with API secret key.
	signature,_ := openssl.HmacSha256Signature(queryParams, API_SECRET)

	// append signature to the query.
	queryParams = queryParams + "&signature=" + signature

	// add API_KEY to the headers.
	var headers = map[string][]string{
		"X-MBX-APIKEY": {API_KEY},
	}

	// parse the query string to url map.
	aQuery,err := url.ParseQuery(queryParams)
	if err != nil {
		logrus.Errorf("Failed to parse the uri %s", err)
		return nil, err
	}

	// perform the query.
	resp, err := cli.get(ctx, "/api/v3/order", aQuery, headers)
	defer ensureReaderClosed(resp)
	if err != nil {
		logrus.Errorf("failed with status code %d", resp.statusCode)
		return nil, err
	}
	err = json.NewDecoder(resp.body).Decode(&orders)
	if err != nil {
		logrus.Errorf("failed to decode json with err %s", err)
	}
	return &orders,nil
}

