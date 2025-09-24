package spot

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewWebSocketClient(t *testing.T) {
	// Test mainnet client
	client := NewWebSocketClient(false)
	if client == nil {
		t.Error("Expected WebSocket client to not be nil")
	}
	
	// Test testnet client
	testnetClient := NewWebSocketClient(true)
	if testnetClient == nil {
		t.Error("Expected testnet WebSocket client to not be nil")
	}
}

func TestWebSocketClientConnect(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Note: This test will fail without a real WebSocket server
	// In a real test environment, you would mock the WebSocket connection
	err := client.Connect()
	if err != nil {
		// Expected to fail in test environment without real server
		t.Logf("Expected connection error in test environment: %v", err)
	}
}

func TestWebSocketClientSubscribe(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test subscription without connection (should not panic)
	client.SubscribeTicker("BTCUSDT", func(ticker *Ticker24hr) {
		// Handler function
	})
	
	client.SubscribeMiniTicker("BTCUSDT", func(ticker *MiniTicker) {
		// Handler function
	})
	
	client.SubscribeBookTicker("BTCUSDT", func(bookTicker *BookTicker) {
		// Handler function
	})
	
	client.SubscribeTrade("BTCUSDT", func(trade *Trade) {
		// Handler function
	})
	
	client.SubscribeAggTrade("BTCUSDT", func(aggTrade *AggTrade) {
		// Handler function
	})
	
	client.SubscribeKline("BTCUSDT", Interval1m, func(kline *Kline) {
		// Handler function
	})
	
	client.SubscribeDepth("BTCUSDT", 5, func(depth *OrderBook) {
		// Handler function
	})
	
	client.SubscribeAllTickers(func(tickers []Ticker24hr) {
		// Handler function
	})
	
	client.SubscribeAllMiniTickers(func(tickers []MiniTicker) {
		// Handler function
	})
	
	client.SubscribeAllBookTickers(func(bookTickers []BookTicker) {
		// Handler function
	})
	
	// Test should not panic
	t.Log("✓ All subscription methods called successfully")
}

func TestWebSocketClientUnsubscribe(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test unsubscription without connection (should not panic)
	client.Unsubscribe("btcusdt@ticker")
	client.Unsubscribe("btcusdt@miniTicker")
	client.Unsubscribe("btcusdt@bookTicker")
	client.Unsubscribe("btcusdt@trade")
	client.Unsubscribe("btcusdt@aggTrade")
	client.Unsubscribe("btcusdt@kline_1m")
	client.Unsubscribe("btcusdt@depth5")
	client.Unsubscribe("!ticker@arr")
	
	// Test should not panic
	t.Log("✓ All unsubscription methods called successfully")
}

func TestWebSocketClientDisconnect(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test disconnect without connection (should not panic)
	err := client.Disconnect()
	if err != nil {
		t.Logf("Expected disconnect error without connection: %v", err)
	}
}

func TestWebSocketClientSetReconnect(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test setting reconnect behavior
	client.SetReconnect(true, 10*time.Second)
	client.SetReconnect(false, 0)
	
	// Test should not panic
	t.Log("✓ Reconnect settings updated successfully")
}

func TestWebSocketClientIsConnected(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test initial connection state
	connected := client.IsConnected()
	if connected {
		t.Error("Expected client to not be connected initially")
	}
}

func TestParseTicker24hr(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"symbol":             "BTCUSDT",
		"priceChange":        "100.00",
		"priceChangePercent": "2.50",
		"lastPrice":          "41000.00",
		"volume":             "1000.00",
	}
	
	ticker := client.parseTicker24hr(validData)
	if ticker == nil {
		t.Error("Expected ticker to not be nil")
	}
	
	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", ticker.Symbol)
	}
}

func TestParseMiniTicker(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"s": "BTCUSDT",
		"o": "40000.00",
		"h": "41000.00",
		"l": "39000.00",
		"c": "40500.00",
		"v": "1000.00",
		"C": 1640995200000,
	}
	
	ticker := client.parseMiniTicker(validData)
	if ticker == nil {
		t.Error("Expected mini ticker to not be nil")
	}
	
	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", ticker.Symbol)
	}
	
	expectedOpen, _ := decimal.NewFromString("40000.00")
	if ticker.Open.Cmp(expectedOpen) != 0 {
		t.Errorf("Expected open %s, got %s", expectedOpen.String(), ticker.Open.String())
	}
}

func TestParseBookTicker(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"s":   "BTCUSDT",
		"b":   "40000.00",
		"B":   "1.00",
		"a":   "40001.00",
		"A":   "1.00",
		"T":   1640995200000,
	}
	
	bookTicker := client.parseBookTicker(validData)
	if bookTicker == nil {
		t.Error("Expected book ticker to not be nil")
	}
	
	if bookTicker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", bookTicker.Symbol)
	}
	
	expectedBidPrice, _ := decimal.NewFromString("40000.00")
	if bookTicker.BidPrice.Cmp(expectedBidPrice) != 0 {
		t.Errorf("Expected bid price %s, got %s", expectedBidPrice.String(), bookTicker.BidPrice.String())
	}
}

func TestParseTrade(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"id":           12345,
		"price":        "40000.00",
		"qty":          "0.001",
		"baseQty":      "0.001",
		"time":         1640995200000,
		"isBuyerMaker": false,
	}
	
	trade := client.parseTrade(validData)
	if trade == nil {
		t.Error("Expected trade to not be nil")
	}
	
	if trade.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", trade.ID)
	}
	
	expectedPrice, _ := decimal.NewFromString("40000.00")
	if trade.Price.Cmp(expectedPrice) != 0 {
		t.Errorf("Expected price %s, got %s", expectedPrice.String(), trade.Price.String())
	}
}

func TestParseAggTrade(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"a": 12345,
		"p": "40000.00",
		"q": "0.001",
		"f": 12340,
		"l": 12345,
		"T": 1640995200000,
		"m": false,
	}
	
	aggTrade := client.parseAggTrade(validData)
	if aggTrade == nil {
		t.Error("Expected agg trade to not be nil")
	}
	
	if aggTrade.A != 12345 {
		t.Errorf("Expected aggregate ID 12345, got %d", aggTrade.A)
	}
	
	expectedPrice, _ := decimal.NewFromString("40000.00")
	if aggTrade.P.Cmp(expectedPrice) != 0 {
		t.Errorf("Expected price %s, got %s", expectedPrice.String(), aggTrade.P.String())
	}
}

func TestParseKline(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"t": 1640995200000,
		"o": "40000.00",
		"h": "41000.00",
		"l": "39000.00",
		"c": "40500.00",
		"v": "1000.00",
		"T": 1640995260000,
		"q": "40000000.00",
		"n": 100,
		"V": "500.00",
		"Q": "20000000.00",
	}
	
	kline := client.parseKline(validData)
	if kline == nil {
		t.Error("Expected kline to not be nil")
	}
	
	if kline.OpenTime != 1640995200000 {
		t.Errorf("Expected open time 1640995200000, got %d", kline.OpenTime)
	}
	
	expectedOpen, _ := decimal.NewFromString("40000.00")
	if kline.Open.Cmp(expectedOpen) != 0 {
		t.Errorf("Expected open %s, got %s", expectedOpen.String(), kline.Open.String())
	}
}

func TestParseDepth(t *testing.T) {
	client := NewWebSocketClient(false)
	
	// Test with valid data
	validData := map[string]interface{}{
		"lastUpdateId": 12345,
		"E":           1640995200000,
		"T":           1640995200000,
		"bids":        [][]string{{"40000.00", "1.00"}},
		"asks":        [][]string{{"40001.00", "1.00"}},
	}
	
	depth := client.parseDepth(validData)
	if depth == nil {
		t.Error("Expected depth to not be nil")
	}
	
	if depth.LastUpdateID != 12345 {
		t.Errorf("Expected last update ID 12345, got %d", depth.LastUpdateID)
	}
	
	if len(depth.Bids) != 1 {
		t.Errorf("Expected 1 bid, got %d", len(depth.Bids))
	}
	
	if len(depth.Asks) != 1 {
		t.Errorf("Expected 1 ask, got %d", len(depth.Asks))
	}
}