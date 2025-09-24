package futures

import (
	"fmt"

	"github.com/asterdex/aster-sdk-go/common"
	"github.com/shopspring/decimal"
)

// Client represents the futures trading client
type Client struct {
	*common.Client
}

// NewClient creates a new futures trading client
func NewClient(config *common.ClientConfig) *Client {
	if config == nil {
		config = common.DefaultConfig()
		config.BaseURL = "https://fapi.asterdex.com"
	}
	
	return &Client{
		Client: common.NewClient(config),
	}
}

// SetTestnet sets the client to use testnet
func (c *Client) SetTestnet(testnet bool) {
	if testnet {
		c.SetBaseURL("https://testnet-fapi.asterdex.com")
	} else {
		c.SetBaseURL("https://fapi.asterdex.com")
	}
}

// Market Data API

// Ping tests connectivity to the REST API
func (c *Client) Ping() error {
	return c.Do("GET", "/fapi/v1/ping", nil, nil, false)
}

// GetServerTime gets the server time
func (c *Client) GetServerTime() (int64, error) {
	var result struct {
		ServerTime int64 `json:"serverTime"`
	}
	err := c.Do("GET", "/fapi/v1/time", nil, &result, false)
	return result.ServerTime, err
}

// GetExchangeInfo gets exchange trading rules and symbol information
func (c *Client) GetExchangeInfo() (*ExchangeInfo, error) {
	var result ExchangeInfo
	err := c.Do("GET", "/fapi/v1/exchangeInfo", nil, &result, false)
	return &result, err
}

// GetOrderBook gets the order book for a symbol
func (c *Client) GetOrderBook(symbol string, limit int) (*OrderBook, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result OrderBook
	err := c.Do("GET", "/fapi/v1/depth", params, &result, false)
	return &result, err
}

// GetRecentTrades gets recent trades for a symbol
func (c *Client) GetRecentTrades(symbol string, limit int) ([]Trade, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []Trade
	err := c.Do("GET", "/fapi/v1/trades", params, &result, false)
	return result, err
}

// GetHistoricalTrades gets historical trades for a symbol
func (c *Client) GetHistoricalTrades(symbol string, limit int, fromID int64) ([]Trade, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	if limit > 0 {
		params["limit"] = limit
	}
	if fromID > 0 {
		params["fromId"] = fromID
	}
	
	var result []Trade
	err := c.Do("GET", "/fapi/v1/historicalTrades", params, &result, true)
	return result, err
}

// GetAggTrades gets compressed/aggregate trades for a symbol
func (c *Client) GetAggTrades(symbol string, fromID, startTime, endTime int64, limit int) ([]AggTrade, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	if fromID > 0 {
		params["fromId"] = fromID
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []AggTrade
	err := c.Do("GET", "/fapi/v1/aggTrades", params, &result, false)
	return result, err
}

// GetKlines gets kline/candlestick data for a symbol
func (c *Client) GetKlines(symbol string, interval KlineInterval, startTime, endTime int64, limit int) ([]Kline, error) {
	params := map[string]any{
		"symbol":   symbol,
		"interval": interval,
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result [][]any
	err := c.Do("GET", "/fapi/v1/klines", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	klines := make([]Kline, len(result))
	for i, k := range result {
		if len(k) < 12 {
			continue
		}
		
		open, _ := decimal.NewFromString(k[1].(string))
		high, _ := decimal.NewFromString(k[2].(string))
		low, _ := decimal.NewFromString(k[3].(string))
		close, _ := decimal.NewFromString(k[4].(string))
		volume, _ := decimal.NewFromString(k[5].(string))
		quoteAssetVolume, _ := decimal.NewFromString(k[7].(string))
		takerBuyBaseAssetVolume, _ := decimal.NewFromString(k[9].(string))
		takerBuyQuoteAssetVolume, _ := decimal.NewFromString(k[10].(string))
		
		klines[i] = Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    close,
			Volume:                   volume,
			CloseTime:                int64(k[6].(float64)),
			QuoteAssetVolume:         quoteAssetVolume,
			NumberOfTrades:           int(k[8].(float64)),
			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolume,
			TakerBuyQuoteAssetVolume: takerBuyQuoteAssetVolume,
		}
	}
	
	return klines, nil
}

// GetIndexPriceKlines gets index price kline/candlestick data
func (c *Client) GetIndexPriceKlines(pair string, interval KlineInterval, startTime, endTime int64, limit int) ([]Kline, error) {
	params := map[string]any{
		"pair":     pair,
		"interval": interval,
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result [][]any
	err := c.Do("GET", "/fapi/v1/indexPriceKlines", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	klines := make([]Kline, len(result))
	for i, k := range result {
		if len(k) < 12 {
			continue
		}
		
		open, _ := decimal.NewFromString(k[1].(string))
		high, _ := decimal.NewFromString(k[2].(string))
		low, _ := decimal.NewFromString(k[3].(string))
		close, _ := decimal.NewFromString(k[4].(string))
		volume, _ := decimal.NewFromString(k[5].(string))
		quoteAssetVolume, _ := decimal.NewFromString(k[7].(string))
		takerBuyBaseAssetVolume, _ := decimal.NewFromString(k[9].(string))
		takerBuyQuoteAssetVolume, _ := decimal.NewFromString(k[10].(string))
		
		klines[i] = Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    close,
			Volume:                   volume,
			CloseTime:                int64(k[6].(float64)),
			QuoteAssetVolume:         quoteAssetVolume,
			NumberOfTrades:           int(k[8].(float64)),
			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolume,
			TakerBuyQuoteAssetVolume: takerBuyQuoteAssetVolume,
		}
	}
	
	return klines, nil
}

// GetMarkPriceKlines gets mark price kline/candlestick data
func (c *Client) GetMarkPriceKlines(symbol string, interval KlineInterval, startTime, endTime int64, limit int) ([]Kline, error) {
	params := map[string]any{
		"symbol":   symbol,
		"interval": interval,
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result [][]any
	err := c.Do("GET", "/fapi/v1/markPriceKlines", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	klines := make([]Kline, len(result))
	for i, k := range result {
		if len(k) < 12 {
			continue
		}
		
		open, _ := decimal.NewFromString(k[1].(string))
		high, _ := decimal.NewFromString(k[2].(string))
		low, _ := decimal.NewFromString(k[3].(string))
		close, _ := decimal.NewFromString(k[4].(string))
		volume, _ := decimal.NewFromString(k[5].(string))
		quoteAssetVolume, _ := decimal.NewFromString(k[7].(string))
		takerBuyBaseAssetVolume, _ := decimal.NewFromString(k[9].(string))
		takerBuyQuoteAssetVolume, _ := decimal.NewFromString(k[10].(string))
		
		klines[i] = Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    close,
			Volume:                   volume,
			CloseTime:                int64(k[6].(float64)),
			QuoteAssetVolume:         quoteAssetVolume,
			NumberOfTrades:           int(k[8].(float64)),
			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolume,
			TakerBuyQuoteAssetVolume: takerBuyQuoteAssetVolume,
		}
	}
	
	return klines, nil
}

// GetMarkPrice gets mark price
func (c *Client) GetMarkPrice(symbol string) (*MarkPrice, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/fapi/v1/premiumIndex", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	// Handle both single mark price and array of mark prices
	switch v := result.(type) {
	case map[string]any:
		var markPrice MarkPrice
		err = c.parseMarkPrice(v, &markPrice)
		return &markPrice, err
	case []any:
		if len(v) > 0 {
			var markPrice MarkPrice
			err = c.parseMarkPrice(v[0].(map[string]any), &markPrice)
			return &markPrice, err
		}
	}
	
	return nil, fmt.Errorf("unexpected response format")
}

// GetAllMarkPrices gets all mark prices
func (c *Client) GetAllMarkPrices() ([]MarkPrice, error) {
	var result []MarkPrice
	err := c.Do("GET", "/fapi/v1/premiumIndex", nil, &result, false)
	return result, err
}

// GetFundingRateHistory gets funding rate history
func (c *Client) GetFundingRateHistory(symbol string, startTime, endTime int64, limit int) ([]FundingRate, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []FundingRate
	err := c.Do("GET", "/fapi/v1/fundingRate", params, &result, false)
	return result, err
}

// GetFundingRateConfig gets funding rate configuration
func (c *Client) GetFundingRateConfig(symbol string) (*FundingRateConfig, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result FundingRateConfig
	err := c.Do("GET", "/fapi/v1/fundingRateConfig", params, &result, false)
	return &result, err
}

// GetTicker24hr gets 24hr ticker price change statistics
func (c *Client) GetTicker24hr(symbol string) (*Ticker24hr, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/fapi/v1/ticker/24hr", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	// Handle both single ticker and array of tickers
	switch v := result.(type) {
	case map[string]any:
		var ticker Ticker24hr
		err = c.parseTicker24hr(v, &ticker)
		return &ticker, err
	case []any:
		if len(v) > 0 {
			var ticker Ticker24hr
			err = c.parseTicker24hr(v[0].(map[string]any), &ticker)
			return &ticker, err
		}
	}
	
	return nil, fmt.Errorf("unexpected response format")
}

// GetAllTickers24hr gets 24hr ticker price change statistics for all symbols
func (c *Client) GetAllTickers24hr() ([]Ticker24hr, error) {
	var result []Ticker24hr
	err := c.Do("GET", "/fapi/v1/ticker/24hr", nil, &result, false)
	return result, err
}

// GetPrice gets the latest price for a symbol
func (c *Client) GetPrice(symbol string) (*PriceTicker, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/fapi/v1/ticker/price", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	// Handle both single price and array of prices
	switch v := result.(type) {
	case map[string]any:
		var price PriceTicker
		err = c.parsePriceTicker(v, &price)
		return &price, err
	case []any:
		if len(v) > 0 {
			var price PriceTicker
			err = c.parsePriceTicker(v[0].(map[string]any), &price)
			return &price, err
		}
	}
	
	return nil, fmt.Errorf("unexpected response format")
}

// GetAllPrices gets the latest price for all symbols
func (c *Client) GetAllPrices() ([]PriceTicker, error) {
	var result []PriceTicker
	err := c.Do("GET", "/fapi/v1/ticker/price", nil, &result, false)
	return result, err
}

// GetBookTicker gets the best bid/ask for a symbol
func (c *Client) GetBookTicker(symbol string) (*BookTicker, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/fapi/v1/ticker/bookTicker", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	// Handle both single book ticker and array of book tickers
	switch v := result.(type) {
	case map[string]any:
		var bookTicker BookTicker
		err = c.parseBookTicker(v, &bookTicker)
		return &bookTicker, err
	case []any:
		if len(v) > 0 {
			var bookTicker BookTicker
			err = c.parseBookTicker(v[0].(map[string]any), &bookTicker)
			return &bookTicker, err
		}
	}
	
	return nil, fmt.Errorf("unexpected response format")
}

// GetAllBookTickers gets the best bid/ask for all symbols
func (c *Client) GetAllBookTickers() ([]BookTicker, error) {
	var result []BookTicker
	err := c.Do("GET", "/fapi/v1/ticker/bookTicker", nil, &result, false)
	return result, err
}

// Account/Trades API

// ChangePositionMode changes position mode
func (c *Client) ChangePositionMode(req *ChangePositionModeRequest) error {
	params := map[string]any{
		"dualSidePosition": req.DualSidePosition,
	}
	return c.Do("POST", "/fapi/v1/positionSide/dual", params, nil, true)
}

// GetCurrentPositionMode gets current position mode
func (c *Client) GetCurrentPositionMode() (bool, error) {
	var result struct {
		DualSidePosition bool `json:"dualSidePosition"`
	}
	err := c.Do("GET", "/fapi/v1/positionSide/dual", nil, &result, true)
	return result.DualSidePosition, err
}

// ChangeMultiAssetsMode changes multi-assets mode
func (c *Client) ChangeMultiAssetsMode(req *ChangeMultiAssetsModeRequest) error {
	params := map[string]any{
		"multiAssetsMargin": req.MultiAssetsMargin,
	}
	return c.Do("POST", "/fapi/v1/multiAssetsMargin", params, nil, true)
}

// GetCurrentMultiAssetsMode gets current multi-assets mode
func (c *Client) GetCurrentMultiAssetsMode() (bool, error) {
	var result struct {
		MultiAssetsMargin bool `json:"multiAssetsMargin"`
	}
	err := c.Do("GET", "/fapi/v1/multiAssetsMargin", nil, &result, true)
	return result.MultiAssetsMargin, err
}

// NewOrder places a new order
func (c *Client) NewOrder(req *NewOrderRequest) (*Order, error) {
	params := map[string]any{
		"symbol": req.Symbol,
		"side":   req.Side,
		"type":   req.Type,
	}
	
	if req.Quantity != "" {
		params["quantity"] = req.Quantity
	}
	if req.QuoteOrderQty != "" {
		params["quoteOrderQty"] = req.QuoteOrderQty
	}
	if req.Price != "" {
		params["price"] = req.Price
	}
	if req.TimeInForce != "" {
		params["timeInForce"] = req.TimeInForce
	}
	if req.NewClientOrderID != "" {
		params["newClientOrderId"] = req.NewClientOrderID
	}
	if req.StopPrice != "" {
		params["stopPrice"] = req.StopPrice
	}
	if req.WorkingType != "" {
		params["workingType"] = req.WorkingType
	}
	if req.PriceProtect {
		params["priceProtect"] = req.PriceProtect
	}
	if req.NewOrderRespType != "" {
		params["newOrderRespType"] = req.NewOrderRespType
	}
	if req.ClosePosition {
		params["closePosition"] = req.ClosePosition
	}
	if req.ReduceOnly {
		params["reduceOnly"] = req.ReduceOnly
	}
	if req.PositionSide != "" {
		params["positionSide"] = req.PositionSide
	}
	
	var result Order
	err := c.Do("POST", "/fapi/v1/order", params, &result, true)
	return &result, err
}

// PlaceMultipleOrders places multiple orders
func (c *Client) PlaceMultipleOrders(orders []NewOrderRequest) ([]Order, error) {
	params := map[string]any{
		"batchOrders": orders,
	}
	
	var result []Order
	err := c.Do("POST", "/fapi/v1/batchOrders", params, &result, true)
	return result, err
}

// Transfer transfers between futures and spot
func (c *Client) Transfer(req *TransferRequest) (*TransferResponse, error) {
	params := map[string]any{
		"asset":  req.Asset,
		"amount": req.Amount,
		"type":   req.Type,
	}
	
	var result TransferResponse
	err := c.Do("POST", "/fapi/v1/transfer", params, &result, true)
	return &result, err
}

// GetOrder gets order information
func (c *Client) GetOrder(symbol string, orderID int64, origClientOrderID string) (*Order, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	
	if orderID > 0 {
		params["orderId"] = orderID
	}
	if origClientOrderID != "" {
		params["origClientOrderId"] = origClientOrderID
	}
	
	var result Order
	err := c.Do("GET", "/fapi/v1/order", params, &result, true)
	return &result, err
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(symbol string, orderID int64, origClientOrderID string) (*Order, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	
	if orderID > 0 {
		params["orderId"] = orderID
	}
	if origClientOrderID != "" {
		params["origClientOrderId"] = origClientOrderID
	}
	
	var result Order
	err := c.Do("DELETE", "/fapi/v1/order", params, &result, true)
	return &result, err
}

// CancelAllOpenOrders cancels all open orders
func (c *Client) CancelAllOpenOrders(symbol string) error {
	params := map[string]any{
		"symbol": symbol,
	}
	return c.Do("DELETE", "/fapi/v1/allOpenOrders", params, nil, true)
}

// CancelMultipleOrders cancels multiple orders
func (c *Client) CancelMultipleOrders(symbol string, orderIDList []int64, origClientOrderIDList []string) ([]Order, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	
	if len(orderIDList) > 0 {
		params["orderIdList"] = orderIDList
	}
	if len(origClientOrderIDList) > 0 {
		params["origClientOrderIdList"] = origClientOrderIDList
	}
	
	var result []Order
	err := c.Do("DELETE", "/fapi/v1/batchOrders", params, &result, true)
	return result, err
}

// AutoCancelAllOpenOrders cancels all open orders with countdown
func (c *Client) AutoCancelAllOpenOrders(symbol string, countdownTime int64) error {
	params := map[string]any{
		"symbol":         symbol,
		"countdownTime": countdownTime,
	}
	return c.Do("POST", "/fapi/v1/countdownCancelAll", params, nil, true)
}

// GetOpenOrders gets all open orders
func (c *Client) GetOpenOrders(symbol string) ([]Order, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result []Order
	err := c.Do("GET", "/fapi/v1/openOrders", params, &result, true)
	return result, err
}

// GetAllOrders gets all orders
func (c *Client) GetAllOrders(symbol string, orderID, startTime, endTime int64, limit int) ([]Order, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	
	if orderID > 0 {
		params["orderId"] = orderID
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []Order
	err := c.Do("GET", "/fapi/v1/allOrders", params, &result, true)
	return result, err
}

// GetAccount gets account information
func (c *Client) GetAccount() (*Account, error) {
	var result Account
	err := c.Do("GET", "/fapi/v2/account", nil, &result, true)
	return &result, err
}

// GetBalance gets futures account balance
func (c *Client) GetBalance() ([]Asset, error) {
	var result []Asset
	err := c.Do("GET", "/fapi/v2/balance", nil, &result, true)
	return result, err
}

// ChangeLeverage changes initial leverage
func (c *Client) ChangeLeverage(req *ChangeLeverageRequest) (*Order, error) {
	params := map[string]any{
		"symbol":   req.Symbol,
		"leverage": req.Leverage,
	}
	
	var result Order
	err := c.Do("POST", "/fapi/v1/leverage", params, &result, true)
	return &result, err
}

// ChangeMarginType changes margin type
func (c *Client) ChangeMarginType(req *ChangeMarginTypeRequest) error {
	params := map[string]any{
		"symbol":     req.Symbol,
		"marginType": req.MarginType,
	}
	return c.Do("POST", "/fapi/v1/marginType", params, nil, true)
}

// ModifyIsolatedPositionMargin modifies isolated position margin
func (c *Client) ModifyIsolatedPositionMargin(req *ModifyIsolatedPositionMarginRequest) (*Order, error) {
	params := map[string]any{
		"symbol":  req.Symbol,
		"amount":  req.Amount,
		"type":    req.Type,
	}
	
	var result Order
	err := c.Do("POST", "/fapi/v1/positionMargin", params, &result, true)
	return &result, err
}

// GetPositionMarginChangeHistory gets position margin change history
func (c *Client) GetPositionMarginChangeHistory(symbol string, startTime, endTime int64, limit int) ([]PositionMarginChangeHistory, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []PositionMarginChangeHistory
	err := c.Do("GET", "/fapi/v1/positionMargin/history", params, &result, true)
	return result, err
}

// GetPositionInfo gets position information
func (c *Client) GetPositionInfo(symbol string) ([]Position, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result []Position
	err := c.Do("GET", "/fapi/v2/positionRisk", params, &result, true)
	return result, err
}

// GetUserTrades gets user trade history
func (c *Client) GetUserTrades(symbol string, orderID, startTime, endTime, fromID int64, limit int) ([]UserTrade, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if orderID > 0 {
		params["orderId"] = orderID
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if fromID > 0 {
		params["fromId"] = fromID
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []UserTrade
	err := c.Do("GET", "/fapi/v1/userTrades", params, &result, true)
	return result, err
}

// GetIncomeHistory gets income history
func (c *Client) GetIncomeHistory(symbol string, incomeType string, startTime, endTime int64, limit int) ([]Income, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if incomeType != "" {
		params["incomeType"] = incomeType
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []Income
	err := c.Do("GET", "/fapi/v1/income", params, &result, true)
	return result, err
}

// GetNotionalBracket gets notional and leverage brackets
func (c *Client) GetNotionalBracket(symbol string) ([]NotionalBracket, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result []NotionalBracket
	err := c.Do("GET", "/fapi/v1/leverageBracket", params, &result, true)
	return result, err
}

// GetADLQuantile gets position ADL quantile estimation
func (c *Client) GetADLQuantile(symbol string) ([]ADLQuantile, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result []ADLQuantile
	err := c.Do("GET", "/fapi/v1/adlQuantile", params, &result, true)
	return result, err
}

// GetForceOrders gets user's force orders
func (c *Client) GetForceOrders(symbol string, autoCloseType string, startTime, endTime int64, limit int) ([]ForceOrder, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if autoCloseType != "" {
		params["autoCloseType"] = autoCloseType
	}
	if startTime > 0 {
		params["startTime"] = startTime
	}
	if endTime > 0 {
		params["endTime"] = endTime
	}
	if limit > 0 {
		params["limit"] = limit
	}
	
	var result []ForceOrder
	err := c.Do("GET", "/fapi/v1/forceOrders", params, &result, true)
	return result, err
}

// GetCommissionRate gets user commission rate
func (c *Client) GetCommissionRate(symbol string) (*CommissionRate, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	
	var result CommissionRate
	err := c.Do("GET", "/fapi/v1/commissionRate", params, &result, true)
	return &result, err
}

// User Data Stream

// CreateListenKey creates a listen key for user data stream
func (c *Client) CreateListenKey() (*ListenKeyResponse, error) {
	var result ListenKeyResponse
	err := c.Do("POST", "/fapi/v1/listenKey", nil, &result, true)
	return &result, err
}

// KeepAliveListenKey keeps the listen key alive
func (c *Client) KeepAliveListenKey(listenKey string) error {
	params := map[string]any{
		"listenKey": listenKey,
	}
	return c.Do("PUT", "/fapi/v1/listenKey", params, nil, true)
}

// CloseListenKey closes the listen key
func (c *Client) CloseListenKey(listenKey string) error {
	params := map[string]any{
		"listenKey": listenKey,
	}
	return c.Do("DELETE", "/fapi/v1/listenKey", params, nil, true)
}

// Helper methods for parsing responses

func (c *Client) parseMarkPrice(data map[string]any, markPrice *MarkPrice) error {
	if symbol, ok := data["symbol"].(string); ok {
		markPrice.Symbol = symbol
	}
	if markPriceStr, ok := data["markPrice"].(string); ok {
		markPrice.MarkPrice, _ = decimal.NewFromString(markPriceStr)
	}
	if indexPrice, ok := data["indexPrice"].(string); ok {
		markPrice.IndexPrice, _ = decimal.NewFromString(indexPrice)
	}
	if estimatedSettlePrice, ok := data["estimatedSettlePrice"].(string); ok {
		markPrice.EstimatedSettlePrice, _ = decimal.NewFromString(estimatedSettlePrice)
	}
	if lastFundingRate, ok := data["lastFundingRate"].(string); ok {
		markPrice.LastFundingRate, _ = decimal.NewFromString(lastFundingRate)
	}
	if nextFundingTime, ok := data["nextFundingTime"].(float64); ok {
		markPrice.NextFundingTime = int64(nextFundingTime)
	}
	if interestRate, ok := data["interestRate"].(string); ok {
		markPrice.InterestRate, _ = decimal.NewFromString(interestRate)
	}
	if timeVal, ok := data["time"].(float64); ok {
		markPrice.Time = int64(timeVal)
	}
	return nil
}

func (c *Client) parseTicker24hr(data map[string]any, ticker *Ticker24hr) error {
	// Parse the ticker data from the map
	// This is a simplified version - in practice you'd want more robust parsing
	if symbol, ok := data["symbol"].(string); ok {
		ticker.Symbol = symbol
	}
	if priceChange, ok := data["priceChange"].(string); ok {
		ticker.PriceChange, _ = decimal.NewFromString(priceChange)
	}
	// Add more fields as needed
	return nil
}

func (c *Client) parsePriceTicker(data map[string]any, price *PriceTicker) error {
	if symbol, ok := data["symbol"].(string); ok {
		price.Symbol = symbol
	}
	if priceStr, ok := data["price"].(string); ok {
		price.Price, _ = decimal.NewFromString(priceStr)
	}
	if timeVal, ok := data["time"].(float64); ok {
		price.Time = int64(timeVal)
	}
	return nil
}

func (c *Client) parseBookTicker(data map[string]any, bookTicker *BookTicker) error {
	if symbol, ok := data["symbol"].(string); ok {
		bookTicker.Symbol = symbol
	}
	if bidPrice, ok := data["bidPrice"].(string); ok {
		bookTicker.BidPrice = bidPrice
	}
	if bidQty, ok := data["bidQty"].(string); ok {
		bookTicker.BidQty = bidQty
	}
	if askPrice, ok := data["askPrice"].(string); ok {
		bookTicker.AskPrice = askPrice
	}
	if askQty, ok := data["askQty"].(string); ok {
		bookTicker.AskQty = askQty
	}
	if timeVal, ok := data["time"].(float64); ok {
		bookTicker.Time = int64(timeVal)
	}
	return nil
}