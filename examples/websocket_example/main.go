package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yiplee/aster-go/futures"
	"github.com/yiplee/aster-go/spot"
)

func main() {
	fmt.Println("=== Aster SDK WebSocket Examples ===")

	// Test spot WebSocket
	testSpotWebSocket()

	// Test futures WebSocket
	testFuturesWebSocket()
}

func testSpotWebSocket() {
	fmt.Println("--- Spot WebSocket Example ---")

	// Create spot WebSocket client
	wsClient := spot.NewWebSocketClient(false) // false for mainnet

	// Connect to WebSocket
	err := wsClient.Connect()
	if err != nil {
		log.Fatal("Failed to connect to spot WebSocket:", err)
	}
	defer wsClient.Disconnect()

	fmt.Println("✓ Connected to spot WebSocket")

	// Subscribe to BTCUSDT ticker
	wsClient.SubscribeTicker("BTCUSDT", func(ticker *spot.Ticker24hr) {
		fmt.Printf("BTCUSDT Ticker: Price=%s, Change=%s%%\n",
			ticker.LastPrice.String(), ticker.PriceChangePercent.String())
	})

	// Subscribe to BTCUSDT mini ticker
	wsClient.SubscribeMiniTicker("BTCUSDT", func(ticker *spot.MiniTicker) {
		fmt.Printf("BTCUSDT Mini Ticker: Open=%s, High=%s, Low=%s, Close=%s, Volume=%s\n",
			ticker.Open.String(), ticker.High.String(), ticker.Low.String(),
			ticker.Close.String(), ticker.Volume.String())
	})

	// Subscribe to BTCUSDT book ticker
	wsClient.SubscribeBookTicker("BTCUSDT", func(bookTicker *spot.BookTicker) {
		fmt.Printf("BTCUSDT Book Ticker: Bid=%s@%s, Ask=%s@%s\n",
			bookTicker.BidQty.String(), bookTicker.BidPrice.String(),
			bookTicker.AskQty.String(), bookTicker.AskPrice.String())
	})

	// Subscribe to BTCUSDT trades
	wsClient.SubscribeTrade("BTCUSDT", func(trade *spot.Trade) {
		side := "SELL"
		if trade.IsBuyerMaker {
			side = "BUY"
		}
		fmt.Printf("BTCUSDT Trade: %s %s @ %s (ID: %d)\n",
			side, trade.Qty.String(), trade.Price.String(), trade.ID)
	})

	// Subscribe to BTCUSDT aggregated trades
	wsClient.SubscribeAggTrade("BTCUSDT", func(aggTrade *spot.AggTrade) {
		side := "SELL"
		if aggTrade.M {
			side = "BUY"
		}
		fmt.Printf("BTCUSDT Agg Trade: %s %s @ %s (ID: %d)\n",
			side, aggTrade.Q.String(), aggTrade.P.String(), aggTrade.A)
	})

	// Subscribe to BTCUSDT klines (1 minute)
	wsClient.SubscribeKline("BTCUSDT", spot.Interval1m, func(kline *spot.Kline) {
		fmt.Printf("BTCUSDT Kline: O=%s H=%s L=%s C=%s V=%s\n",
			kline.Open.String(), kline.High.String(), kline.Low.String(),
			kline.Close.String(), kline.Volume.String())
	})

	// Subscribe to BTCUSDT depth
	wsClient.SubscribeDepth("BTCUSDT", 5, func(depth *spot.OrderBook) {
		fmt.Printf("BTCUSDT Depth: LastUpdateID=%d, Bids=%d, Asks=%d\n",
			depth.LastUpdateID, len(depth.Bids), len(depth.Asks))
	})

	// Subscribe to all tickers
	wsClient.SubscribeAllTickers(func(tickers []spot.Ticker24hr) {
		fmt.Printf("All Tickers: Received %d tickers\n", len(tickers))
		if len(tickers) > 0 {
			fmt.Printf("  First ticker: %s = %s\n", tickers[0].Symbol, tickers[0].LastPrice.String())
		}
	})

	fmt.Println("✓ Subscribed to spot WebSocket streams")
	fmt.Println("  - BTCUSDT ticker")
	fmt.Println("  - BTCUSDT mini ticker")
	fmt.Println("  - BTCUSDT book ticker")
	fmt.Println("  - BTCUSDT trades")
	fmt.Println("  - BTCUSDT aggregated trades")
	fmt.Println("  - BTCUSDT klines (1m)")
	fmt.Println("  - BTCUSDT depth")
	fmt.Println("  - All tickers")

	// Wait for signals
	waitForSignal()
}

func testFuturesWebSocket() {
	fmt.Println("\n--- Futures WebSocket Example ---")

	// Create futures WebSocket client
	wsClient := futures.NewWebSocketClient(false) // false for mainnet

	// Connect to WebSocket
	err := wsClient.Connect()
	if err != nil {
		log.Fatal("Failed to connect to futures WebSocket:", err)
	}
	defer wsClient.Disconnect()

	fmt.Println("✓ Connected to futures WebSocket")

	// Subscribe to BTCUSDT ticker
	wsClient.SubscribeTicker("BTCUSDT", func(ticker *futures.Ticker24hr) {
		fmt.Printf("BTCUSDT Futures Ticker: Price=%s, Change=%s%%\n",
			ticker.LastPrice.String(), ticker.PriceChangePercent.String())
	})

	// Subscribe to BTCUSDT mini ticker
	wsClient.SubscribeMiniTicker("BTCUSDT", func(ticker *futures.MiniTicker) {
		fmt.Printf("BTCUSDT Futures Mini Ticker: Open=%s, High=%s, Low=%s, Close=%s, Volume=%s\n",
			ticker.Open.String(), ticker.High.String(), ticker.Low.String(),
			ticker.Close.String(), ticker.Volume.String())
	})

	// Subscribe to BTCUSDT mark price
	wsClient.SubscribeMarkPrice("BTCUSDT", func(markPrice *futures.MarkPrice) {
		fmt.Printf("BTCUSDT Mark Price: %s (Index: %s, Funding: %s)\n",
			markPrice.MarkPrice.String(), markPrice.IndexPrice.String(), markPrice.LastFundingRate.String())
	})

	// Subscribe to BTCUSDT trades
	wsClient.SubscribeTrade("BTCUSDT", func(trade *futures.Trade) {
		side := "SELL"
		if trade.IsBuyerMaker {
			side = "BUY"
		}
		fmt.Printf("BTCUSDT Futures Trade: %s %s @ %s (ID: %d)\n",
			side, trade.Qty.String(), trade.Price.String(), trade.ID)
	})

	// Subscribe to BTCUSDT aggregated trades
	wsClient.SubscribeAggTrade("BTCUSDT", func(aggTrade *futures.AggTrade) {
		side := "SELL"
		if aggTrade.IsBuyerMaker {
			side = "BUY"
		}
		fmt.Printf("BTCUSDT Futures Agg Trade: %s %s @ %s (ID: %d)\n",
			side, aggTrade.Quantity.String(), aggTrade.Price.String(), aggTrade.AggregateTradeID)
	})

	// Subscribe to BTCUSDT klines (1 minute)
	wsClient.SubscribeKline("BTCUSDT", futures.Interval1m, func(kline *futures.Kline) {
		fmt.Printf("BTCUSDT Futures Kline: O=%s H=%s L=%s C=%s V=%s\n",
			kline.Open.String(), kline.High.String(), kline.Low.String(),
			kline.Close.String(), kline.Volume.String())
	})

	// Subscribe to BTCUSDT depth
	wsClient.SubscribeDepth("BTCUSDT", 5, func(depth *futures.OrderBook) {
		fmt.Printf("BTCUSDT Futures Depth: LastUpdateID=%d, Bids=%d, Asks=%d\n",
			depth.LastUpdateID, len(depth.Bids), len(depth.Asks))
	})

	// Subscribe to all tickers
	wsClient.SubscribeAllTickers(func(tickers []futures.Ticker24hr) {
		fmt.Printf("All Futures Tickers: Received %d tickers\n", len(tickers))
		if len(tickers) > 0 {
			fmt.Printf("  First ticker: %s = %s\n", tickers[0].Symbol, tickers[0].LastPrice.String())
		}
	})

	// Subscribe to all mark prices
	wsClient.SubscribeAllMarkPrices(func(markPrices []futures.MarkPrice) {
		fmt.Printf("All Mark Prices: Received %d mark prices\n", len(markPrices))
		if len(markPrices) > 0 {
			fmt.Printf("  First mark price: %s = %s\n", markPrices[0].Symbol, markPrices[0].MarkPrice.String())
		}
	})

	fmt.Println("✓ Subscribed to futures WebSocket streams")
	fmt.Println("  - BTCUSDT ticker")
	fmt.Println("  - BTCUSDT mini ticker")
	fmt.Println("  - BTCUSDT mark price")
	fmt.Println("  - BTCUSDT trades")
	fmt.Println("  - BTCUSDT aggregated trades")
	fmt.Println("  - BTCUSDT klines (1m)")
	fmt.Println("  - BTCUSDT depth")
	fmt.Println("  - All tickers")
	fmt.Println("  - All mark prices")

	// Wait for signals
	waitForSignal()
}

func waitForSignal() {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-sigChan
	fmt.Println("\n✓ Received interrupt signal, shutting down...")
}
