package futures

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestParseTicker24hr(t *testing.T) {
	// Test with valid data (using WebSocket single letter field names)
	validData := map[string]any{
		"s": "BTCUSDT",
		"P": "100.00",
		"p": "2.50",
		"c": "41000.00",
		"v": "1000.00",
	}

	jsonData, _ := json.Marshal(validData)
	ticker := parseTicker24hr(jsonData)
	if ticker == nil {
		t.Error("Expected ticker to not be nil")
	}

	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", ticker.Symbol)
	}

	expectedPriceChange, _ := decimal.NewFromString("100.00")
	if ticker.PriceChange.Cmp(expectedPriceChange) != 0 {
		t.Errorf("Expected price change %s, got %s", expectedPriceChange.String(), ticker.PriceChange.String())
	}
}

func TestParseBookTicker(t *testing.T) {
	// Test with valid data (using WebSocket single letter field names)
	validData := map[string]any{
		"s": "BTCUSDT",
		"b": "40000.00",
		"B": "1.00",
		"a": "40001.00",
		"A": "1.00",
		"T": 1640995200000,
	}

	jsonData, _ := json.Marshal(validData)
	bookTicker := parseBookTicker(jsonData)
	if bookTicker == nil {
		t.Error("Expected book ticker to not be nil")
	}

	if bookTicker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", bookTicker.Symbol)
	}

	if bookTicker.BidPrice != "40000.00" {
		t.Errorf("Expected bid price 40000.00, got %s", bookTicker.BidPrice)
	}
}

func TestParseTrade(t *testing.T) {
	// Test with valid data (using WebSocket single letter field names)
	validData := map[string]any{
		"t": 12345,
		"p": "40000.00",
		"q": "0.001",
		"b": "0.001",
		"T": 1640995200000,
		"m": false,
	}

	jsonData, _ := json.Marshal(validData)
	trade := parseTrade(jsonData)
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
	// Test with valid data (using WebSocket single letter field names)
	validData := map[string]any{
		"a": 12345,
		"p": "40000.00",
		"q": "0.001",
		"f": 100,
		"l": 200,
		"T": 1640995200000,
		"m": false,
	}

	jsonData, _ := json.Marshal(validData)
	aggTrade := parseAggTrade(jsonData)
	if aggTrade == nil {
		t.Error("Expected agg trade to not be nil")
	}

	if aggTrade.AggregateTradeID != 12345 {
		t.Errorf("Expected aggregate trade ID 12345, got %d", aggTrade.AggregateTradeID)
	}

	expectedPrice, _ := decimal.NewFromString("40000.00")
	if aggTrade.Price.Cmp(expectedPrice) != 0 {
		t.Errorf("Expected price %s, got %s", expectedPrice.String(), aggTrade.Price.String())
	}
}

func TestParseKline(t *testing.T) {
	// Test with valid data (using WebSocket single letter field names)
	validData := map[string]any{
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

	jsonData, _ := json.Marshal(validData)
	kline := parseKline(jsonData)
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
