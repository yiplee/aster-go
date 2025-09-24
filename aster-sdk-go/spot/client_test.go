package spot

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/asterdex/aster-sdk-go/common"
	"github.com/shopspring/decimal"
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
	if config.BaseURL != "https://sapi.asterdex.com" {
		t.Errorf("Expected base URL to be 'https://sapi.asterdex.com', got '%s'", config.BaseURL)
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
	if client.GetConfig().BaseURL != "https://testnet-sapi.asterdex.com" {
		t.Errorf("Expected testnet URL, got %s", client.GetConfig().BaseURL)
	}
	
	// Test setting mainnet
	client.SetTestnet(false)
	if client.GetConfig().BaseURL != "https://sapi.asterdex.com" {
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

func TestGetPrice(t *testing.T) {
	client := NewClient(nil)
	
	// Create a mock response
	responseBody := `{
		"symbol": "BTCUSDT",
		"price": "50000.00",
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
	
	price, err := client.GetPrice("BTCUSDT")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if price.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", price.Symbol)
	}
	
	expectedPrice, _ := decimal.NewFromString("50000.00")
	if price.Price.Cmp(expectedPrice) != 0 {
		t.Errorf("Expected price %s, got %s", expectedPrice.String(), price.Price.String())
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
		"origType": "LIMIT"
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
		Symbol:      "BTCUSDT",
		Side:        OrderSideBuy,
		Type:        OrderTypeLimit,
		Quantity:    decimal.NewFromFloat(1.0),
		Price:       decimal.NewFromFloat(50000.0),
		TimeInForce: TimeInForceGTC,
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
		"balances": [
			{
				"asset": "BTC",
				"free": "4723846.89208129",
				"locked": "0.00000000"
			}
		]
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
	
	if len(account.Balances) != 1 {
		t.Errorf("Expected 1 balance, got %d", len(account.Balances))
	}
	
	expectedFree, _ := decimal.NewFromString("4723846.89208129")
	if account.Balances[0].Free.Cmp(expectedFree) != 0 {
		t.Errorf("Expected free balance %s, got %s", expectedFree.String(), account.Balances[0].Free.String())
	}
}

func TestTransfer(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")
	
	// Create a mock response
	responseBody := `{
		"tranId": 21841,
		"status": "SUCCESS"
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
		Amount:       decimal.NewFromFloat(100.0),
		Asset:        "USDT",
		ClientTranID: "test-tran-id",
		KindType:     "SPOT_FUTURE",
	}
	
	transfer, err := client.Transfer(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if transfer.TranID != 21841 {
		t.Errorf("Expected tran ID 21841, got %d", transfer.TranID)
	}
	
	if transfer.Status != "SUCCESS" {
		t.Errorf("Expected status 'SUCCESS', got '%s'", transfer.Status)
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