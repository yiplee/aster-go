package spot

import "time"

// OrderSide represents the order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType represents the order type
type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeStop   OrderType = "STOP"
	OrderTypeStopMarket OrderType = "STOP_MARKET"
	OrderTypeTakeProfit OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket OrderType = "TAKE_PROFIT_MARKET"
)

// OrderStatus represents the order status
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// TimeInForce represents the time in force
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Canceled
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
	TimeInForceGTX TimeInForce = "GTX" // Good till crossing, Post only
)

// KlineInterval represents the kline interval
type KlineInterval string

const (
	Interval1m  KlineInterval = "1m"
	Interval3m  KlineInterval = "3m"
	Interval5m  KlineInterval = "5m"
	Interval15m KlineInterval = "15m"
	Interval30m KlineInterval = "30m"
	Interval1h  KlineInterval = "1h"
	Interval2h  KlineInterval = "2h"
	Interval4h  KlineInterval = "4h"
	Interval6h  KlineInterval = "6h"
	Interval8h  KlineInterval = "8h"
	Interval12h KlineInterval = "12h"
	Interval1d  KlineInterval = "1d"
	Interval3d  KlineInterval = "3d"
	Interval1w  KlineInterval = "1w"
	Interval1M  KlineInterval = "1M"
)

// ExchangeInfo represents the exchange information
type ExchangeInfo struct {
	Timezone   string    `json:"timezone"`
	ServerTime int64     `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Assets          []struct {
		Asset string `json:"asset"`
	} `json:"assets"`
	Symbols []Symbol `json:"symbols"`
}

// Symbol represents a trading symbol
type Symbol struct {
	Symbol                string   `json:"symbol"`
	Status                string   `json:"status"`
	BaseAsset             string   `json:"baseAsset"`
	QuoteAsset            string   `json:"quoteAsset"`
	PricePrecision        int      `json:"pricePrecision"`
	QuantityPrecision     int      `json:"quantityPrecision"`
	BaseAssetPrecision    int      `json:"baseAssetPrecision"`
	QuotePrecision        int      `json:"quotePrecision"`
	Filters               []Filter `json:"filters"`
	OrderTypes            []string `json:"orderTypes"`
	TimeInForce           []string `json:"timeInForce"`
	OcoAllowed            bool     `json:"ocoAllowed"`
}

// Filter represents a symbol filter
type Filter struct {
	FilterType string `json:"filterType"`
	// Price filter
	MinPrice string `json:"minPrice,omitempty"`
	MaxPrice string `json:"maxPrice,omitempty"`
	TickSize string `json:"tickSize,omitempty"`
	// Lot size filter
	StepSize string `json:"stepSize,omitempty"`
	MinQty   string `json:"minQty,omitempty"`
	MaxQty   string `json:"maxQty,omitempty"`
	// Market lot size filter
	// Percent price filter
	MultiplierUp   string `json:"multiplierUp,omitempty"`
	MultiplierDown string `json:"multiplierDown,omitempty"`
	// Min notional filter
	MinNotional string `json:"minNotional,omitempty"`
	// Max notional filter
	MaxNotional string `json:"maxNotional,omitempty"`
	// Max num orders filter
	Limit int `json:"limit,omitempty"`
}

// OrderBook represents the order book
type OrderBook struct {
	LastUpdateID int64     `json:"lastUpdateId"`
	E            int64     `json:"E"` // Message output time
	T            int64     `json:"T"` // Transaction time
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// Trade represents a trade
type Trade struct {
	ID           int64   `json:"id"`
	Price        string  `json:"price"`
	Qty          string  `json:"qty"`
	BaseQty      string  `json:"baseQty"`
	Time         int64   `json:"time"`
	IsBuyerMaker bool    `json:"isBuyerMaker"`
}

// AggTrade represents an aggregated trade
type AggTrade struct {
	A  int64  `json:"a"` // Aggregate trade ID
	P  string `json:"p"` // Price
	Q  string `json:"q"` // Quantity
	F  int64  `json:"f"` // First trade ID
	L  int64  `json:"l"` // Last trade ID
	T  int64  `json:"T"` // Timestamp
	M  bool   `json:"m"` // Was the buyer the maker?
}

// Kline represents a kline/candlestick
type Kline struct {
	OpenTime                 int64   `json:"openTime"`
	Open                     string  `json:"open"`
	High                     string  `json:"high"`
	Low                      string  `json:"low"`
	Close                    string  `json:"close"`
	Volume                   string  `json:"volume"`
	CloseTime                int64   `json:"closeTime"`
	QuoteAssetVolume         string  `json:"quoteAssetVolume"`
	NumberOfTrades           int     `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  string  `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string  `json:"takerBuyQuoteAssetVolume"`
}

// Ticker24hr represents 24hr ticker statistics
type Ticker24hr struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice            string `json:"bidPrice"`
	BidQty              string `json:"bidQty"`
	AskPrice            string `json:"askPrice"`
	AskQty              string `json:"askQty"`
	OpenPrice           string `json:"openPrice"`
	HighPrice           string `json:"highPrice"`
	LowPrice            string `json:"lowPrice"`
	Volume              string `json:"volume"`
	QuoteVolume         string `json:"quoteVolume"`
	OpenTime            int64  `json:"openTime"`
	CloseTime           int64  `json:"closeTime"`
	FirstID             int64  `json:"firstId"`
	LastID              int64  `json:"lastId"`
	Count               int64  `json:"count"`
	BaseAsset           string `json:"baseAsset"`
	QuoteAsset          string `json:"quoteAsset"`
}

// PriceTicker represents a price ticker
type PriceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

// BookTicker represents the best bid/ask
type BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	Time     int64  `json:"time"`
}

// CommissionRate represents commission rates
type CommissionRate struct {
	Symbol                string `json:"symbol"`
	MakerCommissionRate   string `json:"makerCommissionRate"`
	TakerCommissionRate   string `json:"takerCommissionRate"`
}

// Order represents an order
type Order struct {
	Symbol            string      `json:"symbol"`
	OrderID           int64       `json:"orderId"`
	ClientOrderID     string      `json:"clientOrderId"`
	Price             string      `json:"price"`
	OrigQty           string      `json:"origQty"`
	ExecutedQty       string      `json:"executedQty"`
	CumQuote          string      `json:"cumQuote"`
	Status            OrderStatus `json:"status"`
	TimeInForce       TimeInForce `json:"timeInForce"`
	Type              OrderType   `json:"type"`
	Side              OrderSide   `json:"side"`
	StopPrice         string      `json:"stopPrice"`
	IcebergQty        string      `json:"icebergQty"`
	Time              int64       `json:"time"`
	UpdateTime        int64       `json:"updateTime"`
	IsWorking         bool        `json:"isWorking"`
	OrigQuoteOrderQty string      `json:"origQuoteOrderQty"`
	AvgPrice          string      `json:"avgPrice"`
	OrigType          OrderType   `json:"origType"`
}

// Account represents account information
type Account struct {
	FeeTier      int       `json:"feeTier"`
	CanTrade    bool      `json:"canTrade"`
	CanDeposit  bool      `json:"canDeposit"`
	CanWithdraw bool      `json:"canWithdraw"`
	CanBurnAsset bool     `json:"canBurnAsset"`
	UpdateTime  int64     `json:"updateTime"`
	Balances     []Balance `json:"balances"`
}

// Balance represents an account balance
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// UserTrade represents a user trade
type UserTrade struct {
	Symbol           string `json:"symbol"`
	ID               int64  `json:"id"`
	OrderID          int64  `json:"orderId"`
	Side             string `json:"side"`
	Price            string `json:"price"`
	Qty              string `json:"qty"`
	QuoteQty         string `json:"quoteQty"`
	Commission       string `json:"commission"`
	CommissionAsset  string `json:"commissionAsset"`
	Time             int64  `json:"time"`
	CounterpartyID   int64  `json:"counterpartyId"`
	CreateUpdateID   *int64 `json:"createUpdateId"`
	Maker            bool   `json:"maker"`
	Buyer            bool   `json:"buyer"`
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	Amount        string `json:"amount"`
	Asset         string `json:"asset"`
	ClientTranID  string `json:"clientTranId"`
	KindType      string `json:"kindType"` // FUTURE_SPOT or SPOT_FUTURE
}

// TransferResponse represents a transfer response
type TransferResponse struct {
	TranID  int64  `json:"tranId"`
	Status  string `json:"status"`
}

// WithdrawFeeRequest represents a withdraw fee request
type WithdrawFeeRequest struct {
	ChainID string `json:"chainId"`
	Asset    string `json:"asset"`
}

// WithdrawFeeResponse represents a withdraw fee response
type WithdrawFeeResponse struct {
	TokenPrice   float64 `json:"tokenPrice"`
	GasCost      float64 `json:"gasCost"`
	GasUsdValue  float64 `json:"gasUsdValue"`
}

// WithdrawRequest represents a withdraw request
type WithdrawRequest struct {
	ChainID       string `json:"chainId"`
	Asset         string `json:"asset"`
	Amount        string `json:"amount"`
	Fee           string `json:"fee"`
	Receiver      string `json:"receiver"`
	Nonce         string `json:"nonce"`
	UserSignature string `json:"userSignature"`
}

// WithdrawResponse represents a withdraw response
type WithdrawResponse struct {
	WithdrawID string `json:"withdrawId"`
	Hash       string `json:"hash"`
}

// CreateAPIKeyRequest represents a create API key request
type CreateAPIKeyRequest struct {
	Address           string `json:"address"`
	UserOperationType string `json:"userOperationType"`
	Network           string `json:"network,omitempty"`
	UserSignature     string `json:"userSignature"`
	ApikeyIP          string `json:"apikeyIP,omitempty"`
	Desc              string `json:"desc"`
}

// CreateAPIKeyResponse represents a create API key response
type CreateAPIKeyResponse struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

// ListenKeyResponse represents a listen key response
type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}