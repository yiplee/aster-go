package main

import (
	"fmt"
	"log"
	"time"

	"github.com/asterdex/aster-sdk-go/spot"
	"github.com/asterdex/aster-sdk-go/futures"
)

func main() {
	fmt.Println("=== Aster SDK Market Data Examples ===\n")
	
	// Test both spot and futures APIs
	testSpotMarketData()
	testFuturesMarketData()
}

func testSpotMarketData() {
	fmt.Println("--- Spot Market Data ---")
	
	// Create a spot client
	client := spot.NewClient(nil)
	
	// Test connectivity
	fmt.Println("Testing spot connectivity...")
	err := client.Ping()
	if err != nil {
		log.Fatal("Failed to ping spot server:", err)
	}
	fmt.Println("✓ Spot server is reachable")
	
	// Get server time
	serverTime, err := client.GetServerTime()
	if err != nil {
		log.Fatal("Failed to get spot server time:", err)
	}
	fmt.Printf("✓ Spot server time: %d (%s)\n", serverTime, time.Unix(serverTime/1000, 0).Format(time.RFC3339))
	
	// Get exchange info
	exchangeInfo, err := client.GetExchangeInfo()
	if err != nil {
		log.Fatal("Failed to get spot exchange info:", err)
	}
	fmt.Printf("✓ Spot exchange timezone: %s\n", exchangeInfo.Timezone)
	fmt.Printf("✓ Spot symbols available: %d\n", len(exchangeInfo.Symbols))
	
	// Get all 24hr tickers
	fmt.Println("\nGetting all 24hr tickers...")
	allTickers, err := client.GetAllTickers24hr()
	if err != nil {
		log.Fatal("Failed to get all tickers:", err)
	}
	fmt.Printf("✓ Total tickers: %d\n", len(allTickers))
	
	// Show first 5 tickers
	for i, ticker := range allTickers {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s (change: %s%%)\n", 
			ticker.Symbol, ticker.LastPrice.String(), ticker.PriceChangePercent.String())
	}
	
	// Get all prices
	fmt.Println("\nGetting all prices...")
	allPrices, err := client.GetAllPrices()
	if err != nil {
		log.Fatal("Failed to get all prices:", err)
	}
	fmt.Printf("✓ Total prices: %d\n", len(allPrices))
	
	// Show first 5 prices
	for i, price := range allPrices {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s\n", price.Symbol, price.Price.String())
	}
	
	// Get all book tickers
	fmt.Println("\nGetting all book tickers...")
	allBookTickers, err := client.GetAllBookTickers()
	if err != nil {
		log.Fatal("Failed to get all book tickers:", err)
	}
	fmt.Printf("✓ Total book tickers: %d\n", len(allBookTickers))
	
	// Show first 3 book tickers
	for i, bookTicker := range allBookTickers {
		if i >= 3 {
			break
		}
		fmt.Printf("  %s: Bid=%s@%s Ask=%s@%s\n", 
			bookTicker.Symbol, bookTicker.BidQty.String(), bookTicker.BidPrice.String(),
			bookTicker.AskQty.String(), bookTicker.AskPrice.String())
	}
	
	fmt.Println("✓ Spot market data examples completed\n")
}

func testFuturesMarketData() {
	fmt.Println("--- Futures Market Data ---")
	
	// Create a futures client
	client := futures.NewClient(nil)
	
	// Test connectivity
	fmt.Println("Testing futures connectivity...")
	err := client.Ping()
	if err != nil {
		log.Fatal("Failed to ping futures server:", err)
	}
	fmt.Println("✓ Futures server is reachable")
	
	// Get server time
	serverTime, err := client.GetServerTime()
	if err != nil {
		log.Fatal("Failed to get futures server time:", err)
	}
	fmt.Printf("✓ Futures server time: %d (%s)\n", serverTime, time.Unix(serverTime/1000, 0).Format(time.RFC3339))
	
	// Get exchange info
	exchangeInfo, err := client.GetExchangeInfo()
	if err != nil {
		log.Fatal("Failed to get futures exchange info:", err)
	}
	fmt.Printf("✓ Futures exchange timezone: %s\n", exchangeInfo.Timezone)
	fmt.Printf("✓ Futures symbols available: %d\n", len(exchangeInfo.Symbols))
	
	// Get all mark prices
	fmt.Println("\nGetting all mark prices...")
	allMarkPrices, err := client.GetAllMarkPrices()
	if err != nil {
		log.Fatal("Failed to get all mark prices:", err)
	}
	fmt.Printf("✓ Total mark prices: %d\n", len(allMarkPrices))
	
	// Show first 5 mark prices
	for i, markPrice := range allMarkPrices {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s (index: %s)\n", 
			markPrice.Symbol, markPrice.MarkPrice.String(), markPrice.IndexPrice.String())
	}
	
	// Get all 24hr tickers
	fmt.Println("\nGetting all 24hr tickers...")
	allTickers, err := client.GetAllTickers24hr()
	if err != nil {
		log.Fatal("Failed to get all futures tickers:", err)
	}
	fmt.Printf("✓ Total futures tickers: %d\n", len(allTickers))
	
	// Show first 5 tickers
	for i, ticker := range allTickers {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s (change: %s%%)\n", 
			ticker.Symbol, ticker.LastPrice.String(), ticker.PriceChangePercent.String())
	}
	
	// Get all prices
	fmt.Println("\nGetting all futures prices...")
	allPrices, err := client.GetAllPrices()
	if err != nil {
		log.Fatal("Failed to get all futures prices:", err)
	}
	fmt.Printf("✓ Total futures prices: %d\n", len(allPrices))
	
	// Show first 5 prices
	for i, price := range allPrices {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: %s\n", price.Symbol, price.Price.String())
	}
	
	// Get all book tickers
	fmt.Println("\nGetting all futures book tickers...")
	allBookTickers, err := client.GetAllBookTickers()
	if err != nil {
		log.Fatal("Failed to get all futures book tickers:", err)
	}
	fmt.Printf("✓ Total futures book tickers: %d\n", len(allBookTickers))
	
	// Show first 3 book tickers
	for i, bookTicker := range allBookTickers {
		if i >= 3 {
			break
		}
		fmt.Printf("  %s: Bid=%s@%s Ask=%s@%s\n", 
			bookTicker.Symbol, bookTicker.BidQty.String(), bookTicker.BidPrice.String(),
			bookTicker.AskQty.String(), bookTicker.AskPrice.String())
	}
	
	// Get funding rate history for BTCUSDT
	fmt.Println("\nGetting funding rate history for BTCUSDT...")
	fundingRates, err := client.GetFundingRateHistory("BTCUSDT", 0, 0, 10)
	if err != nil {
		log.Fatal("Failed to get funding rate history:", err)
	}
	fmt.Printf("✓ Funding rate history entries: %d\n", len(fundingRates))
	
	// Show first 3 funding rates
	for i, rate := range fundingRates {
		if i >= 3 {
			break
		}
		fmt.Printf("  %s: %s (time: %d)\n", 
			rate.Symbol, rate.FundingRate.String(), rate.FundingTime)
	}
	
	fmt.Println("✓ Futures market data examples completed\n")
}

func testKlineIntervals() {
	fmt.Println("--- Kline Interval Examples ---")
	
	client := spot.NewClient(nil)
	
	// Test different kline intervals
	intervals := []spot.KlineInterval{
		spot.Interval1m,
		spot.Interval5m,
		spot.Interval15m,
		spot.Interval1h,
		spot.Interval4h,
		spot.Interval1d,
	}
	
	for _, interval := range intervals {
		fmt.Printf("Getting %s klines for BTCUSDT...\n", interval)
		klines, err := client.GetKlines("BTCUSDT", interval, 0, 0, 3)
		if err != nil {
			fmt.Printf("  ✗ Failed to get %s klines: %v\n", interval, err)
			continue
		}
		fmt.Printf("  ✓ Got %d %s klines\n", len(klines), interval)
		
		if len(klines) > 0 {
			kline := klines[0]
			fmt.Printf("    Latest: O=%s H=%s L=%s C=%s V=%s\n",
				kline.Open.String(), kline.High.String(), kline.Low.String(),
				kline.Close.String(), kline.Volume.String())
		}
	}
	
	fmt.Println("✓ Kline interval examples completed")
}