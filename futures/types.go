package futures

import (
	"github.com/shopspring/decimal"
)

// OrderSide represents the order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType represents the order type
type OrderType string

const (
	OrderTypeLimit            OrderType = "LIMIT"
	OrderTypeMarket           OrderType = "MARKET"
	OrderTypeStop             OrderType = "STOP"
	OrderTypeStopMarket       OrderType = "STOP_MARKET"
	OrderTypeTakeProfit       OrderType = "TAKE_PROFIT"
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

// PositionSide represents the position side
type PositionSide string

const (
	PositionSideBoth  PositionSide = "BOTH"
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
)

// WorkingType represents the working type
type WorkingType string

const (
	WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	WorkingTypeLastPrice     WorkingType = "LAST_PRICE"
	WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"
)

// MarginType represents the margin type
type MarginType string

const (
	MarginTypeIsolated MarginType = "ISOLATED"
	MarginTypeCrossed  MarginType = "CROSSED"
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
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []any `json:"exchangeFilters"`
	Assets          []struct {
		Asset string `json:"asset"`
	} `json:"assets"`
	Symbols []Symbol `json:"symbols"`
}

// Symbol represents a trading symbol
type Symbol struct {
	Symbol                 string   `json:"symbol"`
	Status                 string   `json:"status"`
	BaseAsset              string   `json:"baseAsset"`
	QuoteAsset             string   `json:"quoteAsset"`
	PricePrecision         int      `json:"pricePrecision"`
	QuantityPrecision      int      `json:"quantityPrecision"`
	BaseAssetPrecision     int      `json:"baseAssetPrecision"`
	QuotePrecision         int      `json:"quotePrecision"`
	Filters                []Filter `json:"filters"`
	OrderTypes             []string `json:"orderTypes"`
	TimeInForce            []string `json:"timeInForce"`
	OcoAllowed             bool     `json:"ocoAllowed"`
	IsSpotTradingAllowed   bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed bool     `json:"isMarginTradingAllowed"`
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
	LastUpdateID int64      `json:"lastUpdateId"`
	E            int64      `json:"E"` // Message output time
	T            int64      `json:"T"` // Transaction time
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// Trade represents a trade
type Trade struct {
	ID           int64           `json:"id"`
	Price        decimal.Decimal `json:"price"`
	Qty          decimal.Decimal `json:"qty"`
	BaseQty      decimal.Decimal `json:"baseQty"`
	Time         int64           `json:"time"`
	IsBuyerMaker bool            `json:"isBuyerMaker"`
}

// AggTrade represents an aggregated trade
type AggTrade struct {
	A int64           `json:"a"` // Aggregate trade ID
	P decimal.Decimal `json:"p"` // Price
	Q decimal.Decimal `json:"q"` // Quantity
	F int64           `json:"f"` // First trade ID
	L int64           `json:"l"` // Last trade ID
	T int64           `json:"T"` // Timestamp
	M bool            `json:"m"` // Was the buyer the maker?
}

// Kline represents a kline/candlestick
type Kline struct {
	OpenTime                 int64           `json:"openTime"`
	Open                     decimal.Decimal `json:"open"`
	High                     decimal.Decimal `json:"high"`
	Low                      decimal.Decimal `json:"low"`
	Close                    decimal.Decimal `json:"close"`
	Volume                   decimal.Decimal `json:"volume"`
	CloseTime                int64           `json:"closeTime"`
	QuoteAssetVolume         decimal.Decimal `json:"quoteAssetVolume"`
	NumberOfTrades           int             `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  decimal.Decimal `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume decimal.Decimal `json:"takerBuyQuoteAssetVolume"`
}

// Ticker24hr represents 24hr ticker statistics
type Ticker24hr struct {
	Symbol             string          `json:"symbol"`
	PriceChange        decimal.Decimal `json:"priceChange"`
	PriceChangePercent decimal.Decimal `json:"priceChangePercent"`
	WeightedAvgPrice   decimal.Decimal `json:"weightedAvgPrice"`
	PrevClosePrice     decimal.Decimal `json:"prevClosePrice"`
	LastPrice          decimal.Decimal `json:"lastPrice"`
	LastQty            decimal.Decimal `json:"lastQty"`
	BidPrice           decimal.Decimal `json:"bidPrice"`
	BidQty             decimal.Decimal `json:"bidQty"`
	AskPrice           decimal.Decimal `json:"askPrice"`
	AskQty             decimal.Decimal `json:"askQty"`
	OpenPrice          decimal.Decimal `json:"openPrice"`
	HighPrice          decimal.Decimal `json:"highPrice"`
	LowPrice           decimal.Decimal `json:"lowPrice"`
	Volume             decimal.Decimal `json:"volume"`
	QuoteVolume        decimal.Decimal `json:"quoteVolume"`
	OpenTime           int64           `json:"openTime"`
	CloseTime          int64           `json:"closeTime"`
	FirstID            int64           `json:"firstId"`
	LastID             int64           `json:"lastId"`
	Count              int64           `json:"count"`
	BaseAsset          string          `json:"baseAsset"`
	QuoteAsset         string          `json:"quoteAsset"`
}

// PriceTicker represents a price ticker
type PriceTicker struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
	Time   int64           `json:"time"`
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

// MarkPrice represents mark price
type MarkPrice struct {
	Symbol               string          `json:"symbol"`
	MarkPrice            decimal.Decimal `json:"markPrice"`
	IndexPrice           decimal.Decimal `json:"indexPrice"`
	EstimatedSettlePrice decimal.Decimal `json:"estimatedSettlePrice"`
	LastFundingRate      decimal.Decimal `json:"lastFundingRate"`
	NextFundingTime      int64           `json:"nextFundingTime"`
	InterestRate         decimal.Decimal `json:"interestRate"`
	Time                 int64           `json:"time"`
}

// FundingRate represents funding rate
type FundingRate struct {
	Symbol      string          `json:"symbol"`
	FundingRate decimal.Decimal `json:"fundingRate"`
	FundingTime int64           `json:"fundingTime"`
}

// FundingRateConfig represents funding rate configuration
type FundingRateConfig struct {
	Symbol      string `json:"symbol"`
	FundingRate string `json:"fundingRate"`
	FundingTime int64  `json:"fundingTime"`
}

// Order represents an order
type Order struct {
	Symbol            string       `json:"symbol"`
	OrderID           int64        `json:"orderId"`
	ClientOrderID     string       `json:"clientOrderId"`
	Price             string       `json:"price"`
	OrigQty           string       `json:"origQty"`
	ExecutedQty       string       `json:"executedQty"`
	CumQuote          string       `json:"cumQuote"`
	Status            OrderStatus  `json:"status"`
	TimeInForce       TimeInForce  `json:"timeInForce"`
	Type              OrderType    `json:"type"`
	Side              OrderSide    `json:"side"`
	StopPrice         string       `json:"stopPrice"`
	IcebergQty        string       `json:"icebergQty"`
	Time              int64        `json:"time"`
	UpdateTime        int64        `json:"updateTime"`
	IsWorking         bool         `json:"isWorking"`
	OrigQuoteOrderQty string       `json:"origQuoteOrderQty"`
	AvgPrice          string       `json:"avgPrice"`
	OrigType          OrderType    `json:"origType"`
	PositionSide      PositionSide `json:"positionSide"`
	ReduceOnly        bool         `json:"reduceOnly"`
	ClosePosition     bool         `json:"closePosition"`
	WorkingType       WorkingType  `json:"workingType"`
	PriceProtect      bool         `json:"priceProtect"`
}

// Account represents account information
type Account struct {
	FeeTier                     int             `json:"feeTier"`
	CanTrade                    bool            `json:"canTrade"`
	CanDeposit                  bool            `json:"canDeposit"`
	CanWithdraw                 bool            `json:"canWithdraw"`
	CanBurnAsset                bool            `json:"canBurnAsset"`
	UpdateTime                  int64           `json:"updateTime"`
	TotalWalletBalance          decimal.Decimal `json:"totalWalletBalance"`
	TotalUnrealizedProfit       string          `json:"totalUnrealizedProfit"`
	TotalMarginBalance          decimal.Decimal `json:"totalMarginBalance"`
	TotalInitialMargin          string          `json:"totalInitialMargin"`
	TotalMaintMargin            string          `json:"totalMaintMargin"`
	TotalPositionInitialMargin  string          `json:"totalPositionInitialMargin"`
	TotalOpenOrderInitialMargin string          `json:"totalOpenOrderInitialMargin"`
	TotalCrossWalletBalance     string          `json:"totalCrossWalletBalance"`
	TotalCrossUnPnl             string          `json:"totalCrossUnPnl"`
	AvailableBalance            decimal.Decimal `json:"availableBalance"`
	MaxWithdrawAmount           string          `json:"maxWithdrawAmount"`
	Assets                      []Asset         `json:"assets"`
	Positions                   []Position      `json:"positions"`
}

// Asset represents an account asset
type Asset struct {
	Asset                  string          `json:"asset"`
	WalletBalance          decimal.Decimal `json:"walletBalance"`
	UnrealizedProfit       string          `json:"unrealizedProfit"`
	MarginBalance          decimal.Decimal `json:"marginBalance"`
	MaintMargin            string          `json:"maintMargin"`
	InitialMargin          string          `json:"initialMargin"`
	PositionInitialMargin  string          `json:"positionInitialMargin"`
	OpenOrderInitialMargin string          `json:"openOrderInitialMargin"`
	CrossWalletBalance     string          `json:"crossWalletBalance"`
	CrossUnPnl             string          `json:"crossUnPnl"`
	AvailableBalance       decimal.Decimal `json:"availableBalance"`
	MaxWithdrawAmount      string          `json:"maxWithdrawAmount"`
}

// Position represents a position
type Position struct {
	Symbol                 string       `json:"symbol"`
	InitialMargin          string       `json:"initialMargin"`
	MaintMargin            string       `json:"maintMargin"`
	UnrealizedProfit       string       `json:"unrealizedProfit"`
	PositionInitialMargin  string       `json:"positionInitialMargin"`
	OpenOrderInitialMargin string       `json:"openOrderInitialMargin"`
	Leverage               string       `json:"leverage"`
	Isolated               bool         `json:"isolated"`
	EntryPrice             string       `json:"entryPrice"`
	MaxNotional            string       `json:"maxNotional"`
	BidNotional            string       `json:"bidNotional"`
	AskNotional            string       `json:"askNotional"`
	PositionSide           PositionSide `json:"positionSide"`
	PositionAmt            string       `json:"positionAmt"`
	UpdateTime             int64        `json:"updateTime"`
}

// UserTrade represents a user trade
type UserTrade struct {
	Symbol          string          `json:"symbol"`
	ID              int64           `json:"id"`
	OrderID         int64           `json:"orderId"`
	Side            string          `json:"side"`
	Price           decimal.Decimal `json:"price"`
	Qty             decimal.Decimal `json:"qty"`
	QuoteQty        string          `json:"quoteQty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
	Time            int64           `json:"time"`
	CounterpartyID  int64           `json:"counterpartyId"`
	CreateUpdateID  *int64          `json:"createUpdateId"`
	Maker           bool            `json:"maker"`
	Buyer           bool            `json:"buyer"`
	PositionSide    PositionSide    `json:"positionSide"`
}

// Income represents income history
type Income struct {
	Symbol     string `json:"symbol"`
	IncomeType string `json:"incomeType"`
	Income     string `json:"income"`
	Asset      string `json:"asset"`
	Info       string `json:"info"`
	Time       int64  `json:"time"`
	TranID     int64  `json:"tranId"`
	TradeID    string `json:"tradeId"`
}

// LeverageBracket represents leverage bracket
type LeverageBracket struct {
	Bracket          int     `json:"bracket"`
	InitialLeverage  int     `json:"initialLeverage"`
	NotionalCap      int64   `json:"notionalCap"`
	NotionalFloor    int64   `json:"notionalFloor"`
	MaintMarginRatio float64 `json:"maintMarginRatio"`
	Cum              int64   `json:"cum"`
}

// NotionalBracket represents notional bracket
type NotionalBracket struct {
	Symbol   string            `json:"symbol"`
	Brackets []LeverageBracket `json:"brackets"`
}

// ADLQuantile represents ADL quantile
type ADLQuantile struct {
	Symbol      string `json:"symbol"`
	AdlQuantile int    `json:"adlQuantile"`
}

// ForceOrder represents a force order
type ForceOrder struct {
	OrderID       int64           `json:"orderId"`
	Symbol        string          `json:"symbol"`
	Status        string          `json:"status"`
	ClientOrderID string          `json:"clientOrderId"`
	Price         decimal.Decimal `json:"price"`
	AvgPrice      string          `json:"avgPrice"`
	OrigQty       string          `json:"origQty"`
	ExecutedQty   string          `json:"executedQty"`
	OrderStatus   string          `json:"orderStatus"`
	TimeInForce   string          `json:"timeInForce"`
	Type          string          `json:"type"`
	Side          string          `json:"side"`
	StopPrice     string          `json:"stopPrice"`
	Time          int64           `json:"time"`
	UpdateTime    int64           `json:"updateTime"`
}

// CommissionRate represents commission rates
type CommissionRate struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"`
	TakerCommissionRate string `json:"takerCommissionRate"`
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	Asset  string          `json:"asset"`
	Amount decimal.Decimal `json:"amount"`
	Type   int             `json:"type"` // 1: from spot to futures, 2: from futures to spot
}

// TransferResponse represents a transfer response
type TransferResponse struct {
	TranID int64 `json:"tranId"`
}

// ChangePositionModeRequest represents a change position mode request
type ChangePositionModeRequest struct {
	DualSidePosition string `json:"dualSidePosition"` // "true" or "false"
}

// ChangeMultiAssetsModeRequest represents a change multi-assets mode request
type ChangeMultiAssetsModeRequest struct {
	MultiAssetsMargin string `json:"multiAssetsMargin"` // "true" or "false"
}

// ChangeLeverageRequest represents a change leverage request
type ChangeLeverageRequest struct {
	Symbol   string `json:"symbol"`
	Leverage int    `json:"leverage"`
}

// ChangeMarginTypeRequest represents a change margin type request
type ChangeMarginTypeRequest struct {
	Symbol     string `json:"symbol"`
	MarginType string `json:"marginType"`
}

// ModifyIsolatedPositionMarginRequest represents a modify isolated position margin request
type ModifyIsolatedPositionMarginRequest struct {
	Symbol string          `json:"symbol"`
	Amount decimal.Decimal `json:"amount"`
	Type   int             `json:"type"` // 1: add margin, 2: reduce margin
}

// PositionMarginChangeHistory represents position margin change history
type PositionMarginChangeHistory struct {
	Amount       decimal.Decimal `json:"amount"`
	Asset        string          `json:"asset"`
	Symbol       string          `json:"symbol"`
	Time         int64           `json:"time"`
	Type         int             `json:"type"`
	PositionSide string          `json:"positionSide"`
}

// NewOrderRequest represents a new order request
type NewOrderRequest struct {
	Symbol           string          `json:"symbol"`
	Side             OrderSide       `json:"side"`
	Type             OrderType       `json:"type"`
	TimeInForce      TimeInForce     `json:"timeInForce,omitempty"`
	Quantity         decimal.Decimal `json:"quantity,omitempty"`
	QuoteOrderQty    decimal.Decimal `json:"quoteOrderQty,omitempty"`
	Price            decimal.Decimal `json:"price,omitempty"`
	NewClientOrderID string          `json:"newClientOrderId,omitempty"`
	StopPrice        decimal.Decimal `json:"stopPrice,omitempty"`
	WorkingType      WorkingType     `json:"workingType,omitempty"`
	PriceProtect     bool            `json:"priceProtect,omitempty"`
	NewOrderRespType string          `json:"newOrderRespType,omitempty"`
	ClosePosition    bool            `json:"closePosition,omitempty"`
	ReduceOnly       bool            `json:"reduceOnly,omitempty"`
	PositionSide     PositionSide    `json:"positionSide,omitempty"`
}

// ListenKeyResponse represents a listen key response
type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}
