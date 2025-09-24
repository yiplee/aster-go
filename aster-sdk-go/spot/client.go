package spot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/asterdex/aster-sdk-go/common"
)

// Client represents the spot trading client
type Client struct {
	*common.Client
}

// NewClient creates a new spot trading client
func NewClient(config *common.ClientConfig) *Client {
	if config == nil {
		config = common.DefaultConfig()
		config.BaseURL = "https://sapi.asterdex.com"
	}
	
	return &Client{
		Client: common.NewClient(config),
	}
}

// SetTestnet sets the client to use testnet
func (c *Client) SetTestnet(testnet bool) {
	if testnet {
		c.SetBaseURL("https://testnet-sapi.asterdex.com")
	} else {
		c.SetBaseURL("https://sapi.asterdex.com")
	}
}

// Market Data API

// Ping tests connectivity to the REST API
func (c *Client) Ping() error {
	return c.Do("GET", "/api/v1/ping", nil, nil, false)
}

// GetServerTime gets the server time
func (c *Client) GetServerTime() (int64, error) {
	var result struct {
		ServerTime int64 `json:"serverTime"`
	}
	err := c.Do("GET", "/api/v1/time", nil, &result, false)
	return result.ServerTime, err
}

// GetExchangeInfo gets exchange trading rules and symbol information
func (c *Client) GetExchangeInfo() (*ExchangeInfo, error) {
	var result ExchangeInfo
	err := c.Do("GET", "/api/v1/exchangeInfo", nil, &result, false)
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
	err := c.Do("GET", "/api/v1/depth", params, &result, false)
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
	err := c.Do("GET", "/api/v1/trades", params, &result, false)
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
	err := c.Do("GET", "/api/v1/historicalTrades", params, &result, true)
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
	err := c.Do("GET", "/api/v1/aggTrades", params, &result, false)
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
	err := c.Do("GET", "/api/v1/klines", params, &result, false)
	if err != nil {
		return nil, err
	}
	
	klines := make([]Kline, len(result))
	for i, k := range result {
		if len(k) < 12 {
			continue
		}
		
		klines[i] = Kline{
			OpenTime:                 int64(k[0].(float64)),
			Open:                     k[1].(string),
			High:                     k[2].(string),
			Low:                      k[3].(string),
			Close:                    k[4].(string),
			Volume:                   k[5].(string),
			CloseTime:                int64(k[6].(float64)),
			QuoteAssetVolume:         k[7].(string),
			NumberOfTrades:           int(k[8].(float64)),
			TakerBuyBaseAssetVolume:  k[9].(string),
			TakerBuyQuoteAssetVolume: k[10].(string),
		}
	}
	
	return klines, nil
}

// GetTicker24hr gets 24hr ticker price change statistics
func (c *Client) GetTicker24hr(symbol string) (*Ticker24hr, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/api/v1/ticker/24hr", params, &result, false)
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
	err := c.Do("GET", "/api/v1/ticker/24hr", nil, &result, false)
	return result, err
}

// GetPrice gets the latest price for a symbol
func (c *Client) GetPrice(symbol string) (*PriceTicker, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/api/v1/ticker/price", params, &result, false)
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
	err := c.Do("GET", "/api/v1/ticker/price", nil, &result, false)
	return result, err
}

// GetBookTicker gets the best bid/ask for a symbol
func (c *Client) GetBookTicker(symbol string) (*BookTicker, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result any
	err := c.Do("GET", "/api/v1/ticker/bookTicker", params, &result, false)
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
	err := c.Do("GET", "/api/v1/ticker/bookTicker", nil, &result, false)
	return result, err
}

// GetCommissionRate gets the commission rate for a symbol
func (c *Client) GetCommissionRate(symbol string) (*CommissionRate, error) {
	params := map[string]any{
		"symbol": symbol,
	}
	
	var result CommissionRate
	err := c.Do("GET", "/api/v1/commissionRate", params, &result, true)
	return &result, err
}

// Trading API

// NewOrder places a new order
func (c *Client) NewOrder(req *NewOrderRequest) (*Order, error) {
	params := map[string]any{
		"symbol":   req.Symbol,
		"side":     req.Side,
		"type":     req.Type,
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
	if req.NewOrderRespType != "" {
		params["newOrderRespType"] = req.NewOrderRespType
	}
	
	var result Order
	err := c.Do("POST", "/api/v1/order", params, &result, true)
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
	err := c.Do("DELETE", "/api/v1/order", params, &result, true)
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
	err := c.Do("GET", "/api/v1/order", params, &result, true)
	return &result, err
}

// GetOpenOrders gets all open orders
func (c *Client) GetOpenOrders(symbol string) ([]Order, error) {
	params := map[string]any{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	
	var result []Order
	err := c.Do("GET", "/api/v1/openOrders", params, &result, true)
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
	err := c.Do("GET", "/api/v1/allOrders", params, &result, true)
	return result, err
}

// Account API

// GetAccount gets account information
func (c *Client) GetAccount() (*Account, error) {
	var result Account
	err := c.Do("GET", "/api/v1/account", nil, &result, true)
	return &result, err
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
	err := c.Do("GET", "/api/v1/userTrades", params, &result, true)
	return result, err
}

// Transfer between spot and futures
func (c *Client) Transfer(req *TransferRequest) (*TransferResponse, error) {
	params := map[string]any{
		"amount":       req.Amount,
		"asset":        req.Asset,
		"clientTranId": req.ClientTranID,
		"kindType":     req.KindType,
	}
	
	var result TransferResponse
	err := c.Do("POST", "/api/v1/asset/wallet/transfer", params, &result, true)
	return &result, err
}

// GetWithdrawFee gets withdraw fee estimation
func (c *Client) GetWithdrawFee(req *WithdrawFeeRequest) (*WithdrawFeeResponse, error) {
	params := map[string]any{
		"chainId": req.ChainID,
		"asset":   req.Asset,
	}
	
	var result WithdrawFeeResponse
	err := c.Do("GET", "/api/v1/aster/withdraw/estimateFee", params, &result, false)
	return &result, err
}

// Withdraw withdraws assets
func (c *Client) Withdraw(req *WithdrawRequest) (*WithdrawResponse, error) {
	params := map[string]any{
		"chainId":       req.ChainID,
		"asset":         req.Asset,
		"amount":        req.Amount,
		"fee":           req.Fee,
		"receiver":      req.Receiver,
		"nonce":         req.Nonce,
		"userSignature": req.UserSignature,
	}
	
	var result WithdrawResponse
	err := c.Do("POST", "/api/v1/aster/user-withdraw", params, &result, true)
	return &result, err
}

// GetNonce gets nonce for API key creation
func (c *Client) GetNonce(address, userOperationType, network string) (int64, error) {
	params := map[string]any{
		"address":           address,
		"userOperationType": userOperationType,
	}
	if network != "" {
		params["network"] = network
	}
	
	var result int64
	err := c.Do("POST", "/api/v1/getNonce", params, &result, false)
	return result, err
}

// CreateAPIKey creates a new API key
func (c *Client) CreateAPIKey(req *CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	params := map[string]any{
		"address":           req.Address,
		"userOperationType": req.UserOperationType,
		"userSignature":    req.UserSignature,
		"desc":             req.Desc,
	}
	if req.Network != "" {
		params["network"] = req.Network
	}
	if req.ApikeyIP != "" {
		params["apikeyIP"] = req.ApikeyIP
	}
	
	var result CreateAPIKeyResponse
	err := c.Do("POST", "/api/v1/createApiKey", params, &result, false)
	return &result, err
}

// User Data Stream

// CreateListenKey creates a listen key for user data stream
func (c *Client) CreateListenKey() (*ListenKeyResponse, error) {
	var result ListenKeyResponse
	err := c.Do("POST", "/api/v1/listenKey", nil, &result, true)
	return &result, err
}

// KeepAliveListenKey keeps the listen key alive
func (c *Client) KeepAliveListenKey(listenKey string) error {
	params := map[string]any{
		"listenKey": listenKey,
	}
	return c.Do("PUT", "/api/v1/listenKey", params, nil, true)
}

// CloseListenKey closes the listen key
func (c *Client) CloseListenKey(listenKey string) error {
	params := map[string]any{
		"listenKey": listenKey,
	}
	return c.Do("DELETE", "/api/v1/listenKey", params, nil, true)
}

// Helper methods for parsing responses

func (c *Client) parseTicker24hr(data map[string]any, ticker *Ticker24hr) error {
	// Parse the ticker data from the map
	// This is a simplified version - in practice you'd want more robust parsing
	if symbol, ok := data["symbol"].(string); ok {
		ticker.Symbol = symbol
	}
	if priceChange, ok := data["priceChange"].(string); ok {
		ticker.PriceChange = priceChange
	}
	// Add more fields as needed
	return nil
}

func (c *Client) parsePriceTicker(data map[string]any, price *PriceTicker) error {
	if symbol, ok := data["symbol"].(string); ok {
		price.Symbol = symbol
	}
	if priceStr, ok := data["price"].(string); ok {
		price.Price = priceStr
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

// NewOrderRequest represents a new order request
type NewOrderRequest struct {
	Symbol            string      `json:"symbol"`
	Side              OrderSide   `json:"side"`
	Type              OrderType   `json:"type"`
	TimeInForce       TimeInForce `json:"timeInForce,omitempty"`
	Quantity          string      `json:"quantity,omitempty"`
	QuoteOrderQty     string      `json:"quoteOrderQty,omitempty"`
	Price             string      `json:"price,omitempty"`
	NewClientOrderID  string      `json:"newClientOrderId,omitempty"`
	StopPrice         string      `json:"stopPrice,omitempty"`
	NewOrderRespType  string      `json:"newOrderRespType,omitempty"`
}