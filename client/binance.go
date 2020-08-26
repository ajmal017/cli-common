package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vietnamz/cli-common/daemon/openssl"
	"os"
)


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

func (cli *Client) getServerTime( ctx context.Context) (int64, error) {
	var severTime int64
	resp, err := cli.get(ctx, "/api/v3/time", nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return 0, errors.New("Failed to get server time")
	}
	err = json.NewDecoder(resp.body).Decode(&severTime)
	return severTime, err
}
func NewBinanceClient(params string, key string ) *BinanceEndPoint {

	c, _ := NewClientWithOpts(WithHost("https://api.binance.com"))
	serverTime, err := c.getServerTime(context.Background())
	if err != nil {
		logrus.Errorf("Failed to get server time %s", err)
		return nil
	}
	if key != "" {
		sign, err := openssl.HmacSha256Signature(params, key)
		if err != nil {
			logrus.Errorf("Failed to convert params to signature with %s", params)
			return nil
		}
		return &BinanceEndPoint{
			Key: key,
			signature: sign,
			ServerTime: serverTime,
		}
		// fall back to os envs.
	} else if binanceSecretKey := os.Getenv("BINANCE_KEY"); binanceSecretKey != "" {

	}
	return nil
}

