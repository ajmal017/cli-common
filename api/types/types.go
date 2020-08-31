package types

// alias to a builder version
// format MAJOR_MINOR_DATE
// ex: 2020_1_2020/08/18
type BuilderVersion string

// PrimeNumber is used to retrieve the prime number.
// GET "/prime?number=123"
// "number" is the upper bound value.
type PrimeNumber struct {
	Number     string       `json:"number"`
	Mtime      string   `json:"elapsed_time"`
}
// Version contains response of Engine API:
// GET "/version"
type Version struct {
	Platform   struct{ Name string } `json:",omitempty"`

	// The following fields are deprecated, they relate to the Engine component and are kept for backwards compatibility

	Version       string
	APIVersion    string `json:"ApiVersion"`
	MinAPIVersion string `json:"MinAPIVersion,omitempty"`
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string `json:",omitempty"`
	Experimental  bool   `json:",omitempty"`
	BuildTime     string `json:",omitempty"`
}

// Ping contains response of Engine API:
// GET "/ping"
type Ping struct {
	// API Version
	// ex: v1
	APIVersion     string
	// OS Type
	// ex: linux
	OSType         string
	// Not release yet.
	Experimental   bool
	//! refs above.
	BuilderVersion BuilderVersion
}

type BinanceServerTime struct {
	ServerTime     int64    `json:"serverTime"`
}
/*
{
  "symbol": "LTCBTC",
  "orderId": 1,
  "orderListId": -1 //Unless part of an OCO, the value will always be -1.
  "clientOrderId": "myOrder1",
  "price": "0.1",
  "origQty": "1.0",
  "executedQty": "0.0",
  "cummulativeQuoteQty": "0.0",
  "status": "NEW",
  "timeInForce": "GTC",
  "type": "LIMIT",
  "side": "BUY",
  "stopPrice": "0.0",
  "icebergQty": "0.0",
  "time": 1499827319559,
  "updateTime": 1499827319559,
  "isWorking": true,
  "origQuoteOrderQty": "0.000000"
}
 */
type BinanceOrder struct {
	Symbol string `json:"symbol"`
	OrderId int64 `json:"orderId"`
	OrderListId int64 `json:"orderListId"`
	ClientOrderId string `json:"clientOrderId"`
	Price string `json:"price"`
	OrigQty string `json:"origQty"`
	ExecutedQty string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status string `json:"status"`
	TimeInForce string `json:"timeInForce"`
	Type string `json:"type"`
	Side string `json:"side"`
	StopPrice string `json:"stopPrice"`
	IcebergQty string `json:"icebergQty"`
	Time int64 `json:"time"`
	UpdateTime int64 `json:"updateTime"`
	IsWorking bool `json:"isWorking"`
	OrigQuoteOrderQty string `json:"origQuoteOrderQty"`
}

type BinanceOrderShort struct {
	Symbol string `json:"symbol"`
	OrderId int64 `json:"orderId"`
	ClientOrderId string `json:"clientOrderId"`
}

type BinanceOrderParams struct {
	Symbol string `url:"symbol"`
	OrderId int64 `url:"orderId,omitempty"`
	OrigClientOrderId string `url:"origClientOrderId,omitempty"`
	RecvWindow int64 `url:"recvWindow,omitempty"`
	Timestamp int64 `url:"timestamp"`
}

type BinanceAllOrdersParams struct {
	Symbol string `url:"symbol"`
	OrderId int64 `url:"orderId,omitempty"`
	StartTime int64 `url:"startTime,omitempty"`
	EndTime int64 `url:"endTime,omitempty"`
	Limit int `url:"limit,omitempty"`
	RecvWindow int64 `url:"recvWindow,omitempty"`
	Timestamp int64	`url:"timestamp"`
}

