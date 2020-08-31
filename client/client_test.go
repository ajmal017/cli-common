package client

import (
	"context"
	"testing"
)

/*
func TestNewClientWithOpts(t *testing.T) {
	cli, err := NewClientWithOpts(WithHost("http://localhost:8080"))
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err = cli.Ping(context.Background())
	if err != nil {
		t.FailNow()
	}
}
*/

func TestNewBinanceClient_Ping(t *testing.T) {
	cli, err := NewClientWithOpts(WithHost("https://api.binance.com"), WithVersion(""))
	if err != nil {
		t.Errorf("%s", err)
	}

	ping := cli.PingBinanceServer(context.Background())
	if ping == false {
		t.Fatalf("Binance server is down.")
	}
}

func TestNewBinanceClient_Time(t *testing.T)  {
	cli, err := NewClientWithOpts(WithHost("https://api.binance.com/"), WithVersion(""))
	if err != nil {
		t.Errorf("%s", err)
	}
	binanceServerTime, err := cli.GetBinanceServerTime(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}
	if binanceServerTime == 0 {
		t.Errorf("Failed to get binance server time")
	}
	t.Logf("Binance server time is %d", binanceServerTime)
}

func TestNewBinanceClient_Get(t *testing.T)  {
	cli, err := NewClientWithOpts(WithHost("https://api.binance.com/"), WithVersion(""))
	if err != nil {
		t.Errorf("%s", err)
	}
	binanceServerTime, err := cli.GetOrder(context.WithCancel())
	if err != nil {
		t.Fatalf(err.Error())
	}
}




/*
func TestNewClientWithOpts(t *testing.T) {
	uri := types.BinanceOrderParams {
		Symbol: "ABC",
		OrigClientOrderId: "",
		OrderId: 0,
	}
	v, _ := query.Values(uri)
	t.Errorf("Failed with %s", v.Encode())
}
*/
