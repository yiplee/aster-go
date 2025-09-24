package futures

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/yiplee/aster-go/common"
)

// WebSocketClient represents the futures WebSocket client
type WebSocketClient struct {
	*common.WebSocketClient
	baseURL string
}

// NewWebSocketClient creates a new futures WebSocket client
func NewWebSocketClient(testnet bool) *WebSocketClient {
	baseURL := "wss://fstream.asterdex.com"
	if testnet {
		baseURL = "wss://testnet-fstream.asterdex.com"
	}

	return &WebSocketClient{
		WebSocketClient: common.NewWebSocketClient(baseURL),
		baseURL:         baseURL,
	}
}

// Subscribe to individual symbol ticker streams
func (c *WebSocketClient) SubscribeTicker(symbol string, handler func(*Ticker24hr)) {
	stream := fmt.Sprintf("%s@ticker", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if ticker := c.parseTicker24hr(data); ticker != nil {
			handler(ticker)
		}
	})
}

// Subscribe to all symbols ticker stream
func (c *WebSocketClient) SubscribeAllTickers(handler func([]Ticker24hr)) {
	stream := "!ticker@arr"
	c.Subscribe(stream, func(data any) {
		if tickers := c.parseAllTickers(data); tickers != nil {
			handler(tickers)
		}
	})
}

// Subscribe to individual symbol mini ticker streams
func (c *WebSocketClient) SubscribeMiniTicker(symbol string, handler func(*MiniTicker)) {
	stream := fmt.Sprintf("%s@miniTicker", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if ticker := c.parseMiniTicker(data); ticker != nil {
			handler(ticker)
		}
	})
}

// Subscribe to all symbols mini ticker stream
func (c *WebSocketClient) SubscribeAllMiniTickers(handler func([]MiniTicker)) {
	stream := "!miniTicker@arr"
	c.Subscribe(stream, func(data any) {
		if tickers := c.parseAllMiniTickers(data); tickers != nil {
			handler(tickers)
		}
	})
}

// Subscribe to individual symbol book ticker streams
func (c *WebSocketClient) SubscribeBookTicker(symbol string, handler func(*BookTicker)) {
	stream := fmt.Sprintf("%s@bookTicker", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if bookTicker := c.parseBookTicker(data); bookTicker != nil {
			handler(bookTicker)
		}
	})
}

// Subscribe to all symbols book ticker stream
func (c *WebSocketClient) SubscribeAllBookTickers(handler func([]BookTicker)) {
	stream := "!bookTicker@arr"
	c.Subscribe(stream, func(data any) {
		if bookTickers := c.parseAllBookTickers(data); bookTickers != nil {
			handler(bookTickers)
		}
	})
}

// Subscribe to individual symbol trade streams
func (c *WebSocketClient) SubscribeTrade(symbol string, handler func(*Trade)) {
	stream := fmt.Sprintf("%s@trade", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if trade := c.parseTrade(data); trade != nil {
			handler(trade)
		}
	})
}

// Subscribe to individual symbol aggregated trade streams
func (c *WebSocketClient) SubscribeAggTrade(symbol string, handler func(*AggTrade)) {
	stream := fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if aggTrade := c.parseAggTrade(data); aggTrade != nil {
			handler(aggTrade)
		}
	})
}

// Subscribe to individual symbol kline streams
func (c *WebSocketClient) SubscribeKline(symbol string, interval KlineInterval, handler func(*Kline)) {
	stream := fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	c.Subscribe(stream, func(data any) {
		if kline := c.parseKline(data); kline != nil {
			handler(kline)
		}
	})
}

// Subscribe to individual symbol depth streams
func (c *WebSocketClient) SubscribeDepth(symbol string, levels int, handler func(*OrderBook)) {
	var stream string
	if levels > 0 {
		stream = fmt.Sprintf("%s@depth%d", strings.ToLower(symbol), levels)
	} else {
		stream = fmt.Sprintf("%s@depth", strings.ToLower(symbol))
	}
	c.Subscribe(stream, func(data any) {
		if depth := c.parseDepth(data); depth != nil {
			handler(depth)
		}
	})
}

// Subscribe to individual symbol depth streams with 100ms updates
func (c *WebSocketClient) SubscribeDepthWithUpdates(symbol string, levels int, handler func(*OrderBook)) {
	var stream string
	if levels > 0 {
		stream = fmt.Sprintf("%s@depth%d@100ms", strings.ToLower(symbol), levels)
	} else {
		stream = fmt.Sprintf("%s@depth@100ms", strings.ToLower(symbol))
	}
	c.Subscribe(stream, func(data any) {
		if depth := c.parseDepth(data); depth != nil {
			handler(depth)
		}
	})
}

// Subscribe to mark price streams
func (c *WebSocketClient) SubscribeMarkPrice(symbol string, handler func(*MarkPrice)) {
	stream := fmt.Sprintf("%s@markPrice", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if markPrice := c.parseMarkPrice(data); markPrice != nil {
			handler(markPrice)
		}
	})
}

// Subscribe to all symbols mark price stream
func (c *WebSocketClient) SubscribeAllMarkPrices(handler func([]MarkPrice)) {
	stream := "!markPrice@arr"
	c.Subscribe(stream, func(data any) {
		if markPrices := c.parseAllMarkPrices(data); markPrices != nil {
			handler(markPrices)
		}
	})
}

// Subscribe to funding rate streams
func (c *WebSocketClient) SubscribeFundingRate(symbol string, handler func(*FundingRate)) {
	stream := fmt.Sprintf("%s@markPrice", strings.ToLower(symbol))
	c.Subscribe(stream, func(data any) {
		if fundingRate := c.parseFundingRate(data); fundingRate != nil {
			handler(fundingRate)
		}
	})
}

// MiniTicker represents a mini ticker
type MiniTicker struct {
	Symbol    string          `json:"s"`
	Open      decimal.Decimal `json:"o"`
	High      decimal.Decimal `json:"h"`
	Low       decimal.Decimal `json:"l"`
	Close     decimal.Decimal `json:"c"`
	Volume    decimal.Decimal `json:"v"`
	CloseTime int64           `json:"C"`
}

// Parse methods for WebSocket data

func (c *WebSocketClient) parseTicker24hr(data any) *Ticker24hr {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var ticker Ticker24hr
	if err := json.Unmarshal(jsonData, &ticker); err != nil {
		return nil
	}

	return &ticker
}

func (c *WebSocketClient) parseAllTickers(data any) []Ticker24hr {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var tickers []Ticker24hr
	if err := json.Unmarshal(jsonData, &tickers); err != nil {
		return nil
	}

	return tickers
}

func (c *WebSocketClient) parseMiniTicker(data any) *MiniTicker {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var ticker MiniTicker
	if err := json.Unmarshal(jsonData, &ticker); err != nil {
		return nil
	}

	return &ticker
}

func (c *WebSocketClient) parseAllMiniTickers(data any) []MiniTicker {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var tickers []MiniTicker
	if err := json.Unmarshal(jsonData, &tickers); err != nil {
		return nil
	}

	return tickers
}

func (c *WebSocketClient) parseBookTicker(data any) *BookTicker {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var bookTicker BookTicker
	if err := json.Unmarshal(jsonData, &bookTicker); err != nil {
		return nil
	}

	return &bookTicker
}

func (c *WebSocketClient) parseAllBookTickers(data any) []BookTicker {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var bookTickers []BookTicker
	if err := json.Unmarshal(jsonData, &bookTickers); err != nil {
		return nil
	}

	return bookTickers
}

func (c *WebSocketClient) parseTrade(data any) *Trade {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var trade Trade
	if err := json.Unmarshal(jsonData, &trade); err != nil {
		return nil
	}

	return &trade
}

func (c *WebSocketClient) parseAggTrade(data any) *AggTrade {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var aggTrade AggTrade
	if err := json.Unmarshal(jsonData, &aggTrade); err != nil {
		return nil
	}

	return &aggTrade
}

func (c *WebSocketClient) parseKline(data any) *Kline {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var kline Kline
	if err := json.Unmarshal(jsonData, &kline); err != nil {
		return nil
	}

	return &kline
}

func (c *WebSocketClient) parseDepth(data any) *OrderBook {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var depth OrderBook
	if err := json.Unmarshal(jsonData, &depth); err != nil {
		return nil
	}

	return &depth
}

func (c *WebSocketClient) parseMarkPrice(data any) *MarkPrice {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var markPrice MarkPrice
	if err := json.Unmarshal(jsonData, &markPrice); err != nil {
		return nil
	}

	return &markPrice
}

func (c *WebSocketClient) parseAllMarkPrices(data any) []MarkPrice {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var markPrices []MarkPrice
	if err := json.Unmarshal(jsonData, &markPrices); err != nil {
		return nil
	}

	return markPrices
}

func (c *WebSocketClient) parseFundingRate(data any) *FundingRate {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var fundingRate FundingRate
	if err := json.Unmarshal(jsonData, &fundingRate); err != nil {
		return nil
	}

	return &fundingRate
}
