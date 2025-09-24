package futures

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/yiplee/aster-go/common"
)

// MockHTTPClient is a mock implementation of HTTPClient
type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func TestNewClient(t *testing.T) {
	client := NewClient(nil)

	if client == nil {
		t.Error("Expected client to not be nil")
	}

	config := client.GetConfig()
	if config.BaseURL != "https://fapi.asterdex.com" {
		t.Errorf("Expected base URL to be 'https://fapi.asterdex.com', got '%s'", config.BaseURL)
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := &common.ClientConfig{
		BaseURL: "https://custom.example.com",
	}

	client := NewClient(config)

	if client == nil {
		t.Error("Expected client to not be nil")
	}

	if client.GetConfig().BaseURL != config.BaseURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", config.BaseURL, client.GetConfig().BaseURL)
	}
}

func TestSetTestnet(t *testing.T) {
	client := NewClient(nil)

	// Test setting testnet
	client.SetTestnet(true)
	if client.GetConfig().BaseURL != "https://testnet-fapi.asterdex.com" {
		t.Errorf("Expected testnet URL, got %s", client.GetConfig().BaseURL)
	}

	// Test setting mainnet
	client.SetTestnet(false)
	if client.GetConfig().BaseURL != "https://fapi.asterdex.com" {
		t.Errorf("Expected mainnet URL, got %s", client.GetConfig().BaseURL)
	}
}

func TestPing(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	err := client.Ping()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetServerTime(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `{"serverTime": 1640995200000}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	serverTime, err := client.GetServerTime()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if serverTime != 1640995200000 {
		t.Errorf("Expected server time 1640995200000, got %d", serverTime)
	}
}

func TestGetExchangeInfo(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `{
		"timezone": "UTC",
		"serverTime": 1640995200000,
		"rateLimits": [],
		"exchangeFilters": [],
		"assets": [],
		"symbols": []
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	exchangeInfo, err := client.GetExchangeInfo()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exchangeInfo.Timezone != "UTC" {
		t.Errorf("Expected timezone 'UTC', got '%s'", exchangeInfo.Timezone)
	}
}

func TestGetOrderBook(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `{
		"lastUpdateId": 1027024,
		"E": 1589436922972,
		"T": 1589436922959,
		"bids": [["4.00000000", "431.00000000"]],
		"asks": [["4.00000200", "12.00000000"]]
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	orderBook, err := client.GetOrderBook("BTCUSDT", 100)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if orderBook.LastUpdateID != 1027024 {
		t.Errorf("Expected last update ID 1027024, got %d", orderBook.LastUpdateID)
	}

	if len(orderBook.Bids) != 1 {
		t.Errorf("Expected 1 bid, got %d", len(orderBook.Bids))
	}
}

func TestGetRecentTrades(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `[
		{
			"id": 657,
			"price": "1.01000000",
			"qty": "5.00000000",
			"baseQty": "4.95049505",
			"time": 1755156533943,
			"isBuyerMaker": false
		}
	]`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	trades, err := client.GetRecentTrades("BTCUSDT", 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(trades) != 1 {
		t.Errorf("Expected 1 trade, got %d", len(trades))
	}

	expectedPrice, _ := decimal.NewFromString("1.01000000")
	if trades[0].Price.Cmp(expectedPrice) != 0 {
		t.Errorf("Expected price %s, got %s", expectedPrice.String(), trades[0].Price.String())
	}
}

func TestGetKlines(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `[
		[
			1499040000000,
			"0.01634790",
			"0.80000000",
			"0.01575800",
			"0.01577100",
			"148976.11427815",
			1499644799999,
			"2434.19055334",
			308,
			"1756.87402397",
			"28.46694368"
		]
	]`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	klines, err := client.GetKlines("BTCUSDT", Interval1h, 0, 0, 100)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(klines) != 1 {
		t.Errorf("Expected 1 kline, got %d", len(klines))
	}

	expectedOpen, _ := decimal.NewFromString("0.01634790")
	if klines[0].Open.Cmp(expectedOpen) != 0 {
		t.Errorf("Expected open %s, got %s", expectedOpen.String(), klines[0].Open.String())
	}
}

func TestGetMarkPrice(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `{
		"symbol": "BTCUSDT",
		"markPrice": "50000.00",
		"indexPrice": "50001.00",
		"estimatedSettlePrice": "50002.00",
		"lastFundingRate": "0.0001",
		"nextFundingTime": 1640995200000,
		"interestRate": "0.0001",
		"time": 1640995200000
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	markPrice, err := client.GetMarkPrice("BTCUSDT")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if markPrice.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", markPrice.Symbol)
	}

	expectedMarkPrice, _ := decimal.NewFromString("50000.00")
	if markPrice.MarkPrice.Cmp(expectedMarkPrice) != 0 {
		t.Errorf("Expected mark price %s, got %s", expectedMarkPrice.String(), markPrice.MarkPrice.String())
	}
}

func TestGetTicker24hr(t *testing.T) {
	client := NewClient(nil)

	// Create a mock response
	responseBody := `{
		"symbol": "BTCUSDT",
		"priceChange": "-94.99999800",
		"priceChangePercent": "-95.960",
		"weightedAvgPrice": "0.29628482",
		"prevClosePrice": "3.89000000",
		"lastPrice": "4.00000200",
		"lastQty": "200.00000000",
		"bidPrice": "866.66000000",
		"bidQty": "72.05100000",
		"askPrice": "866.73000000",
		"askQty": "1.21700000",
		"openPrice": "99.00000000",
		"highPrice": "100.00000000",
		"lowPrice": "0.10000000",
		"volume": "8913.30000000",
		"quoteVolume": "15.30000000",
		"openTime": 1499783499040,
		"closeTime": 1499869899040,
		"firstId": 28385,
		"lastId": 28460,
		"count": 76,
		"baseAsset": "BTC",
		"quoteAsset": "USDT"
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	ticker, err := client.GetTicker24hr("BTCUSDT")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", ticker.Symbol)
	}

	expectedPriceChange, _ := decimal.NewFromString("-94.99999800")
	if ticker.PriceChange.Cmp(expectedPriceChange) != 0 {
		t.Errorf("Expected price change %s, got %s", expectedPriceChange.String(), ticker.PriceChange.String())
	}
}

func TestNewOrder(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")

	// Create a mock response
	responseBody := `{
		"symbol": "BTCUSDT",
		"orderId": 28,
		"clientOrderId": "6gCrw2kRUAF9CvJDGP16IP",
		"price": "50000.00",
		"origQty": "1.0",
		"executedQty": "0.0",
		"cumQuote": "0.0",
		"status": "NEW",
		"timeInForce": "GTC",
		"type": "LIMIT",
		"side": "BUY",
		"stopPrice": "0.0",
		"icebergQty": "0.0",
		"time": 1507725176595,
		"updateTime": 1507725176595,
		"isWorking": true,
		"origQuoteOrderQty": "0.0",
		"avgPrice": "0.0",
		"origType": "LIMIT",
		"positionSide": "BOTH",
		"reduceOnly": false,
		"closePosition": false,
		"workingType": "CONTRACT_PRICE",
		"priceProtect": false
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	req := &NewOrderRequest{
		Symbol:   "BTCUSDT",
		Side:     OrderSideBuy,
		Type:     OrderTypeLimit,
		Quantity: decimal.NewFromFloat(1.0),
		Price:    decimal.NewFromFloat(50000.0),
	}

	order, err := client.NewOrder(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if order.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", order.Symbol)
	}

	if order.OrderID != 28 {
		t.Errorf("Expected order ID 28, got %d", order.OrderID)
	}
}

func TestGetAccount(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")

	// Create a mock response
	responseBody := `{
		"feeTier": 0,
		"canTrade": true,
		"canDeposit": true,
		"canWithdraw": true,
		"canBurnAsset": true,
		"updateTime": 0,
		"totalWalletBalance": "1000.00",
		"totalUnrealizedProfit": "100.00",
		"totalMarginBalance": "1100.00",
		"totalInitialMargin": "100.00",
		"totalMaintMargin": "50.00",
		"totalPositionInitialMargin": "50.00",
		"totalOpenOrderInitialMargin": "50.00",
		"totalCrossWalletBalance": "1000.00",
		"totalCrossUnPnl": "100.00",
		"availableBalance": "1000.00",
		"maxWithdrawAmount": "1000.00",
		"assets": [
			{
				"asset": "USDT",
				"walletBalance": "1000.00",
				"unrealizedProfit": "100.00",
				"marginBalance": "1100.00",
				"maintMargin": "50.00",
				"initialMargin": "100.00",
				"positionInitialMargin": "50.00",
				"openOrderInitialMargin": "50.00",
				"crossWalletBalance": "1000.00",
				"crossUnPnl": "100.00",
				"availableBalance": "1000.00",
				"maxWithdrawAmount": "1000.00"
			}
		],
		"positions": []
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	account, err := client.GetAccount()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !account.CanTrade {
		t.Error("Expected canTrade to be true")
	}

	expectedTotalWalletBalance, _ := decimal.NewFromString("1000.00")
	if account.TotalWalletBalance.Cmp(expectedTotalWalletBalance) != 0 {
		t.Errorf("Expected total wallet balance %s, got %s", expectedTotalWalletBalance.String(), account.TotalWalletBalance.String())
	}

	if len(account.Assets) != 1 {
		t.Errorf("Expected 1 asset, got %d", len(account.Assets))
	}
}

func TestChangePositionMode(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")

	// Create a mock response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	req := &ChangePositionModeRequest{
		DualSidePosition: "true",
	}

	err := client.ChangePositionMode(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetCurrentPositionMode(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")

	// Create a mock response
	responseBody := `{
		"dualSidePosition": true
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	dualSide, err := client.GetCurrentPositionMode()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !dualSide {
		t.Error("Expected dual side position to be true")
	}
}

func TestTransfer(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")

	// Create a mock response
	responseBody := `{
		"tranId": 21841
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	req := &TransferRequest{
		Asset:  "USDT",
		Amount: decimal.NewFromFloat(100.0),
		Type:   1, // Spot to futures
	}

	transfer, err := client.Transfer(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if transfer.TranID != 21841 {
		t.Errorf("Expected tran ID 21841, got %d", transfer.TranID)
	}
}

func TestCreateListenKey(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")

	// Create a mock response
	responseBody := `{
		"listenKey": "pqia91ma19a5s61cv6a81va65sdf19v8a65a1a5s61cv6a81va65sdf19v8a65a1"
	}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)

	listenKey, err := client.CreateListenKey()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedKey := "pqia91ma19a5s61cv6a81va65sdf19v8a65a1a5s61cv6a81va65sdf19v8a65a1"
	if listenKey.ListenKey != expectedKey {
		t.Errorf("Expected listen key '%s', got '%s'", expectedKey, listenKey.ListenKey)
	}
}
