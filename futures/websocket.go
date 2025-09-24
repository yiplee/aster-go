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
	c.Subscribe(stream, func(data json.RawMessage) {
		if ticker := parseTicker24hr(data); ticker != nil {
			handler(ticker)
		}
	})
}

// Subscribe to all symbols ticker stream
func (c *WebSocketClient) SubscribeAllTickers(handler func([]Ticker24hr)) {
	stream := "!ticker@arr"
	c.Subscribe(stream, func(data json.RawMessage) {
		if tickers := parseAllTickers(data); tickers != nil {
			handler(tickers)
		}
	})
}

// Subscribe to individual symbol mini ticker streams
func (c *WebSocketClient) SubscribeMiniTicker(symbol string, handler func(*MiniTicker)) {
	stream := fmt.Sprintf("%s@miniTicker", strings.ToLower(symbol))
	c.Subscribe(stream, func(data json.RawMessage) {
		if ticker := parseMiniTicker(data); ticker != nil {
			handler(ticker)
		}
	})
}

// Subscribe to all symbols mini ticker stream
func (c *WebSocketClient) SubscribeAllMiniTickers(handler func([]MiniTicker)) {
	stream := "!miniTicker@arr"
	c.Subscribe(stream, func(data json.RawMessage) {
		if tickers := parseAllMiniTickers(data); tickers != nil {
			handler(tickers)
		}
	})
}

// Subscribe to individual symbol book ticker streams
func (c *WebSocketClient) SubscribeBookTicker(symbol string, handler func(*BookTicker)) {
	stream := fmt.Sprintf("%s@bookTicker", strings.ToLower(symbol))
	c.Subscribe(stream, func(data json.RawMessage) {
		if bookTicker := parseBookTicker(data); bookTicker != nil {
			handler(bookTicker)
		}
	})
}

// Subscribe to all symbols book ticker stream
func (c *WebSocketClient) SubscribeAllBookTickers(handler func([]BookTicker)) {
	stream := "!bookTicker@arr"
	c.Subscribe(stream, func(data json.RawMessage) {
		if bookTickers := parseAllBookTickers(data); bookTickers != nil {
			handler(bookTickers)
		}
	})
}

// Subscribe to individual symbol trade streams
func (c *WebSocketClient) SubscribeTrade(symbol string, handler func(*Trade)) {
	stream := fmt.Sprintf("%s@trade", strings.ToLower(symbol))
	c.Subscribe(stream, func(data json.RawMessage) {
		if trade := parseTrade(data); trade != nil {
			handler(trade)
		}
	})
}

// Subscribe to individual symbol aggregated trade streams
func (c *WebSocketClient) SubscribeAggTrade(symbol string, handler func(*AggTrade)) {
	stream := fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
	c.Subscribe(stream, func(data json.RawMessage) {
		if aggTrade := parseAggTrade(data); aggTrade != nil {
			handler(aggTrade)
		}
	})
}

// Subscribe to individual symbol kline streams
func (c *WebSocketClient) SubscribeKline(symbol string, interval KlineInterval, handler func(*Kline)) {
	stream := fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	c.Subscribe(stream, func(data json.RawMessage) {
		if kline := parseKline(data); kline != nil {
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
	c.Subscribe(stream, func(data json.RawMessage) {
		if depth := parseDepth(data); depth != nil {
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
	c.Subscribe(stream, func(data json.RawMessage) {
		if depth := parseDepth(data); depth != nil {
			handler(depth)
		}
	})
}

// Subscribe to mark price streams
func (c *WebSocketClient) SubscribeMarkPrice(symbol string, handler func(*MarkPrice)) {
	stream := fmt.Sprintf("%s@markPrice", strings.ToLower(symbol))
	c.Subscribe(stream, func(data json.RawMessage) {
		if markPrice := parseMarkPrice(data); markPrice != nil {
			handler(markPrice)
		}
	})
}

// Subscribe to all symbols mark price stream
func (c *WebSocketClient) SubscribeAllMarkPrices(handler func([]MarkPrice)) {
	stream := "!markPrice@arr"
	c.Subscribe(stream, func(data json.RawMessage) {
		if markPrices := parseAllMarkPrices(data); markPrices != nil {
			handler(markPrices)
		}
	})
}

// Subscribe to funding rate streams
func (c *WebSocketClient) SubscribeFundingRate(symbol string, handler func(*FundingRate)) {
	stream := fmt.Sprintf("%s@markPrice", strings.ToLower(symbol))
	c.Subscribe(stream, func(data json.RawMessage) {
		if fundingRate := parseFundingRate(data); fundingRate != nil {
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

// Parse functions for WebSocket data

func parseTicker24hr(data json.RawMessage) *Ticker24hr {
	var rawData map[string]any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil
	}

	ticker := &Ticker24hr{}

	// Parse WebSocket format (single letter fields)
	if symbol, ok := rawData["s"].(string); ok {
		ticker.Symbol = symbol
	}
	if priceChange, ok := rawData["P"].(string); ok {
		ticker.PriceChange, _ = decimal.NewFromString(priceChange)
	}
	if priceChangePercent, ok := rawData["p"].(string); ok {
		ticker.PriceChangePercent, _ = decimal.NewFromString(priceChangePercent)
	}
	if weightedAvgPrice, ok := rawData["w"].(string); ok {
		ticker.WeightedAvgPrice, _ = decimal.NewFromString(weightedAvgPrice)
	}
	if prevClosePrice, ok := rawData["x"].(string); ok {
		ticker.PrevClosePrice, _ = decimal.NewFromString(prevClosePrice)
	}
	if lastPrice, ok := rawData["c"].(string); ok {
		ticker.LastPrice, _ = decimal.NewFromString(lastPrice)
	}
	if lastQty, ok := rawData["Q"].(string); ok {
		ticker.LastQty, _ = decimal.NewFromString(lastQty)
	}
	if bidPrice, ok := rawData["b"].(string); ok {
		ticker.BidPrice, _ = decimal.NewFromString(bidPrice)
	}
	if bidQty, ok := rawData["B"].(string); ok {
		ticker.BidQty, _ = decimal.NewFromString(bidQty)
	}
	if askPrice, ok := rawData["a"].(string); ok {
		ticker.AskPrice, _ = decimal.NewFromString(askPrice)
	}
	if askQty, ok := rawData["A"].(string); ok {
		ticker.AskQty, _ = decimal.NewFromString(askQty)
	}
	if openPrice, ok := rawData["o"].(string); ok {
		ticker.OpenPrice, _ = decimal.NewFromString(openPrice)
	}
	if highPrice, ok := rawData["h"].(string); ok {
		ticker.HighPrice, _ = decimal.NewFromString(highPrice)
	}
	if lowPrice, ok := rawData["l"].(string); ok {
		ticker.LowPrice, _ = decimal.NewFromString(lowPrice)
	}
	if volume, ok := rawData["v"].(string); ok {
		ticker.Volume, _ = decimal.NewFromString(volume)
	}
	if quoteVolume, ok := rawData["q"].(string); ok {
		ticker.QuoteVolume, _ = decimal.NewFromString(quoteVolume)
	}
	if openTime, ok := rawData["O"].(float64); ok {
		ticker.OpenTime = int64(openTime)
	}
	if closeTime, ok := rawData["C"].(float64); ok {
		ticker.CloseTime = int64(closeTime)
	}
	if firstID, ok := rawData["F"].(float64); ok {
		ticker.FirstID = int64(firstID)
	}
	if lastID, ok := rawData["L"].(float64); ok {
		ticker.LastID = int64(lastID)
	}
	if count, ok := rawData["n"].(float64); ok {
		ticker.Count = int64(count)
	}
	if baseAsset, ok := rawData["baseAsset"].(string); ok {
		ticker.BaseAsset = baseAsset
	}
	if quoteAsset, ok := rawData["quoteAsset"].(string); ok {
		ticker.QuoteAsset = quoteAsset
	}

	return ticker
}

func parseAllTickers(data json.RawMessage) []Ticker24hr {
	var tickers []Ticker24hr
	if err := json.Unmarshal(data, &tickers); err != nil {
		return nil
	}

	return tickers
}

func parseMiniTicker(data json.RawMessage) *MiniTicker {
	var ticker MiniTicker
	if err := json.Unmarshal(data, &ticker); err != nil {
		return nil
	}

	return &ticker
}

func parseAllMiniTickers(data json.RawMessage) []MiniTicker {
	var tickers []MiniTicker
	if err := json.Unmarshal(data, &tickers); err != nil {
		return nil
	}

	return tickers
}

func parseBookTicker(data json.RawMessage) *BookTicker {
	var rawData map[string]any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil
	}

	bookTicker := &BookTicker{}

	// Parse WebSocket format (single letter fields)
	if symbol, ok := rawData["s"].(string); ok {
		bookTicker.Symbol = symbol
	}
	if bidPrice, ok := rawData["b"].(string); ok {
		bookTicker.BidPrice = bidPrice
	}
	if bidQty, ok := rawData["B"].(string); ok {
		bookTicker.BidQty = bidQty
	}
	if askPrice, ok := rawData["a"].(string); ok {
		bookTicker.AskPrice = askPrice
	}
	if askQty, ok := rawData["A"].(string); ok {
		bookTicker.AskQty = askQty
	}
	if time, ok := rawData["T"].(float64); ok {
		bookTicker.Time = int64(time)
	}

	return bookTicker
}

func parseAllBookTickers(data json.RawMessage) []BookTicker {
	var bookTickers []BookTicker
	if err := json.Unmarshal(data, &bookTickers); err != nil {
		return nil
	}

	return bookTickers
}

func parseTrade(data json.RawMessage) *Trade {
	var rawData map[string]any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil
	}

	trade := &Trade{}

	// Parse WebSocket format (single letter fields)
	if id, ok := rawData["t"].(float64); ok {
		trade.ID = int64(id)
	}
	if price, ok := rawData["p"].(string); ok {
		trade.Price, _ = decimal.NewFromString(price)
	}
	if qty, ok := rawData["q"].(string); ok {
		trade.Qty, _ = decimal.NewFromString(qty)
	}
	if baseQty, ok := rawData["b"].(string); ok {
		trade.BaseQty, _ = decimal.NewFromString(baseQty)
	}
	if time, ok := rawData["T"].(float64); ok {
		trade.Time = int64(time)
	}
	if isBuyerMaker, ok := rawData["m"].(bool); ok {
		trade.IsBuyerMaker = isBuyerMaker
	}

	return trade
}

func parseAggTrade(data json.RawMessage) *AggTrade {
	var rawData map[string]any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil
	}

	aggTrade := &AggTrade{}

	// Parse WebSocket format (single letter fields)
	if aggregateTradeID, ok := rawData["a"].(float64); ok {
		aggTrade.AggregateTradeID = int64(aggregateTradeID)
	}
	if price, ok := rawData["p"].(string); ok {
		aggTrade.Price, _ = decimal.NewFromString(price)
	}
	if quantity, ok := rawData["q"].(string); ok {
		aggTrade.Quantity, _ = decimal.NewFromString(quantity)
	}
	if firstTradeID, ok := rawData["f"].(float64); ok {
		aggTrade.FirstTradeID = int64(firstTradeID)
	}
	if lastTradeID, ok := rawData["l"].(float64); ok {
		aggTrade.LastTradeID = int64(lastTradeID)
	}
	if timestamp, ok := rawData["T"].(float64); ok {
		aggTrade.Timestamp = int64(timestamp)
	}
	if isBuyerMaker, ok := rawData["m"].(bool); ok {
		aggTrade.IsBuyerMaker = isBuyerMaker
	}

	return aggTrade
}

func parseKline(data json.RawMessage) *Kline {
	var rawData map[string]any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil
	}

	kline := &Kline{}

	// Parse WebSocket format (single letter fields)
	if openTime, ok := rawData["t"].(float64); ok {
		kline.OpenTime = int64(openTime)
	}
	if open, ok := rawData["o"].(string); ok {
		kline.Open, _ = decimal.NewFromString(open)
	}
	if high, ok := rawData["h"].(string); ok {
		kline.High, _ = decimal.NewFromString(high)
	}
	if low, ok := rawData["l"].(string); ok {
		kline.Low, _ = decimal.NewFromString(low)
	}
	if close, ok := rawData["c"].(string); ok {
		kline.Close, _ = decimal.NewFromString(close)
	}
	if volume, ok := rawData["v"].(string); ok {
		kline.Volume, _ = decimal.NewFromString(volume)
	}
	if closeTime, ok := rawData["T"].(float64); ok {
		kline.CloseTime = int64(closeTime)
	}
	if quoteAssetVolume, ok := rawData["q"].(string); ok {
		kline.QuoteAssetVolume, _ = decimal.NewFromString(quoteAssetVolume)
	}
	if numberOfTrades, ok := rawData["n"].(float64); ok {
		kline.NumberOfTrades = int(numberOfTrades)
	}
	if takerBuyBaseAssetVolume, ok := rawData["V"].(string); ok {
		kline.TakerBuyBaseAssetVolume, _ = decimal.NewFromString(takerBuyBaseAssetVolume)
	}
	if takerBuyQuoteAssetVolume, ok := rawData["Q"].(string); ok {
		kline.TakerBuyQuoteAssetVolume, _ = decimal.NewFromString(takerBuyQuoteAssetVolume)
	}

	return kline
}

func parseDepth(data json.RawMessage) *OrderBook {
	var rawData map[string]any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil
	}

	depth := &OrderBook{}

	// Parse WebSocket format (mixed field names)
	if lastUpdateID, ok := rawData["lastUpdateId"].(float64); ok {
		depth.LastUpdateID = int64(lastUpdateID)
	}
	if messageTime, ok := rawData["E"].(float64); ok {
		depth.MessageTime = int64(messageTime)
	}
	if transactionTime, ok := rawData["T"].(float64); ok {
		depth.TransactionTime = int64(transactionTime)
	}
	if bids, ok := rawData["bids"].([]any); ok {
		for _, bid := range bids {
			if bidSlice, ok := bid.([]any); ok && len(bidSlice) >= 2 {
				bidStr := []string{}
				for _, item := range bidSlice {
					if str, ok := item.(string); ok {
						bidStr = append(bidStr, str)
					}
				}
				if len(bidStr) >= 2 {
					depth.Bids = append(depth.Bids, bidStr)
				}
			}
		}
	}
	if asks, ok := rawData["asks"].([]any); ok {
		for _, ask := range asks {
			if askSlice, ok := ask.([]any); ok && len(askSlice) >= 2 {
				askStr := []string{}
				for _, item := range askSlice {
					if str, ok := item.(string); ok {
						askStr = append(askStr, str)
					}
				}
				if len(askStr) >= 2 {
					depth.Asks = append(depth.Asks, askStr)
				}
			}
		}
	}

	return depth
}

func parseMarkPrice(data json.RawMessage) *MarkPrice {
	var markPrice MarkPrice
	if err := json.Unmarshal(data, &markPrice); err != nil {
		return nil
	}

	return &markPrice
}

func parseAllMarkPrices(data json.RawMessage) []MarkPrice {
	var markPrices []MarkPrice
	if err := json.Unmarshal(data, &markPrices); err != nil {
		return nil
	}

	return markPrices
}

func parseFundingRate(data json.RawMessage) *FundingRate {
	var fundingRate FundingRate
	if err := json.Unmarshal(data, &fundingRate); err != nil {
		return nil
	}

	return &fundingRate
}
