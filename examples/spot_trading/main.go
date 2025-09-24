package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yiplee/aster-go/spot"
)

func main() {
	// Create a spot client
	client := spot.NewClient(nil)

	// Set API credentials (replace with your actual credentials)
	client.SetAPIKey("your-api-key", "your-secret-key")

	// Test connectivity
	fmt.Println("Testing connectivity...")
	err := client.Ping()
	if err != nil {
		log.Fatal("Failed to ping server:", err)
	}
	fmt.Println("✓ Server is reachable")

	// Get server time
	fmt.Println("\nGetting server time...")
	serverTime, err := client.GetServerTime()
	if err != nil {
		log.Fatal("Failed to get server time:", err)
	}
	fmt.Printf("✓ Server time: %d (%s)\n", serverTime, time.Unix(serverTime/1000, 0).Format(time.RFC3339))

	// Get exchange info
	fmt.Println("\nGetting exchange information...")
	exchangeInfo, err := client.GetExchangeInfo()
	if err != nil {
		log.Fatal("Failed to get exchange info:", err)
	}
	fmt.Printf("✓ Exchange timezone: %s\n", exchangeInfo.Timezone)
	fmt.Printf("✓ Server time: %d\n", exchangeInfo.ServerTime)
	fmt.Printf("✓ Number of symbols: %d\n", len(exchangeInfo.Symbols))

	// Get 24hr ticker for BTCUSDT
	fmt.Println("\nGetting 24hr ticker for BTCUSDT...")
	ticker, err := client.GetTicker24hr("BTCUSDT")
	if err != nil {
		log.Fatal("Failed to get ticker:", err)
	}
	fmt.Printf("✓ Symbol: %s\n", ticker.Symbol)
	fmt.Printf("✓ Last price: %s\n", ticker.LastPrice.String())
	fmt.Printf("✓ Price change: %s (%s%%)\n", ticker.PriceChange.String(), ticker.PriceChangePercent.String())
	fmt.Printf("✓ Volume: %s\n", ticker.Volume.String())

	// Get order book
	fmt.Println("\nGetting order book for BTCUSDT...")
	orderBook, err := client.GetOrderBook("BTCUSDT", 5)
	if err != nil {
		log.Fatal("Failed to get order book:", err)
	}
	fmt.Printf("✓ Last update ID: %d\n", orderBook.LastUpdateID)
	fmt.Printf("✓ Number of bids: %d\n", len(orderBook.Bids))
	fmt.Printf("✓ Number of asks: %d\n", len(orderBook.Asks))

	if len(orderBook.Bids) > 0 {
		fmt.Printf("✓ Best bid: %s @ %s\n", orderBook.Bids[0][1], orderBook.Bids[0][0])
	}
	if len(orderBook.Asks) > 0 {
		fmt.Printf("✓ Best ask: %s @ %s\n", orderBook.Asks[0][1], orderBook.Asks[0][0])
	}

	// Get recent trades
	fmt.Println("\nGetting recent trades for BTCUSDT...")
	trades, err := client.GetRecentTrades("BTCUSDT", 5)
	if err != nil {
		log.Fatal("Failed to get recent trades:", err)
	}
	fmt.Printf("✓ Number of trades: %d\n", len(trades))

	for i, trade := range trades {
		if i >= 3 { // Show only first 3 trades
			break
		}
		side := "SELL"
		if trade.IsBuyerMaker {
			side = "BUY"
		}
		fmt.Printf("✓ Trade %d: %s %s @ %s (time: %d)\n",
			trade.ID, side, trade.Qty.String(), trade.Price.String(), trade.Time)
	}

	// Get kline data
	fmt.Println("\nGetting kline data for BTCUSDT...")
	klines, err := client.GetKlines("BTCUSDT", spot.Interval1h, 0, 0, 5)
	if err != nil {
		log.Fatal("Failed to get klines:", err)
	}
	fmt.Printf("✓ Number of klines: %d\n", len(klines))

	for i, kline := range klines {
		if i >= 3 { // Show only first 3 klines
			break
		}
		fmt.Printf("✓ Kline %d: O=%s H=%s L=%s C=%s V=%s\n",
			i+1, kline.Open.String(), kline.High.String(),
			kline.Low.String(), kline.Close.String(), kline.Volume.String())
	}

	// Example of placing an order (commented out for safety)
	/*
		fmt.Println("\nPlacing a test order...")
		orderReq := &spot.NewOrderRequest{
			Symbol:      "BTCUSDT",
			Side:        spot.OrderSideBuy,
			Type:        spot.OrderTypeLimit,
			Quantity:    decimal.NewFromFloat(0.001),
			Price:       decimal.NewFromFloat(30000.0), // Low price to avoid execution
			TimeInForce: spot.TimeInForceGTC,
		}

		order, err := client.NewOrder(orderReq)
		if err != nil {
			log.Fatal("Failed to place order:", err)
		}
		fmt.Printf("✓ Order placed: ID=%d, Status=%s\n", order.OrderID, order.Status)
	*/

	fmt.Println("\n✓ All examples completed successfully!")
}
