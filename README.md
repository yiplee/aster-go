# Aster SDK for Go

A comprehensive Go SDK for the Aster Finance exchange, providing both spot and futures trading capabilities with high precision decimal arithmetic.

## Features

- **Spot Trading**: Complete spot trading API implementation
- **Futures Trading**: Full futures/contracts trading support
- **High Precision**: Uses `decimal.Decimal` for all price and quantity calculations
- **Type Safety**: Strongly typed API with comprehensive error handling
- **Comprehensive Testing**: Full test coverage with mock implementations
- **WebSocket Support**: Real-time data streams with automatic reconnection
- **Rate Limiting**: Built-in rate limit handling
- **Authentication**: HMAC SHA256 signature support

## Installation

```bash
go get github.com/yiplee/aster-go
```

## Quick Start

### Spot Trading

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yiplee/aster-go/spot"
    "github.com/shopspring/decimal"
)

func main() {
    // Create a spot client
    client := spot.NewClient(nil)
    
    // Set API credentials
    client.SetAPIKey("your-api-key", "your-secret-key")
    
    // Test connectivity
    err := client.Ping()
    if err != nil {
        log.Fatal(err)
    }
    
    // Get server time
    serverTime, err := client.GetServerTime()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Server time: %d\n", serverTime)
    
    // Get exchange info
    exchangeInfo, err := client.GetExchangeInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Exchange timezone: %s\n", exchangeInfo.Timezone)
    
    // Get 24hr ticker
    ticker, err := client.GetTicker24hr("BTCUSDT")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTCUSDT price: %s\n", ticker.LastPrice.String())
    
    // Place a limit order
    orderReq := &spot.NewOrderRequest{
        Symbol:      "BTCUSDT",
        Side:        spot.OrderSideBuy,
        Type:        spot.OrderTypeLimit,
        Quantity:    decimal.NewFromFloat(0.001),
        Price:       decimal.NewFromFloat(50000.0),
        TimeInForce: spot.TimeInForceGTC,
    }
    
    order, err := client.NewOrder(orderReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order placed: %d\n", order.OrderID)
}
```

### Futures Trading

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yiplee/aster-go/futures"
    "github.com/shopspring/decimal"
)

func main() {
    // Create a futures client
    client := futures.NewClient(nil)
    
    // Set API credentials
    client.SetAPIKey("your-api-key", "your-secret-key")
    
    // Test connectivity
    err := client.Ping()
    if err != nil {
        log.Fatal(err)
    }
    
    // Get mark price
    markPrice, err := client.GetMarkPrice("BTCUSDT")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTCUSDT mark price: %s\n", markPrice.MarkPrice.String())
    
    // Place a futures order
    orderReq := &futures.NewOrderRequest{
        Symbol:   "BTCUSDT",
        Side:     futures.OrderSideBuy,
        Type:     futures.OrderTypeLimit,
        Quantity: decimal.NewFromFloat(0.01),
        Price:    decimal.NewFromFloat(50000.0),
    }
    
    order, err := client.NewOrder(orderReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Futures order placed: %d\n", order.OrderID)
}
```

### WebSocket Real-time Data

```go
package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/yiplee/aster-go/spot"
)

func main() {
    // Create WebSocket client
    wsClient := spot.NewWebSocketClient(false) // false for mainnet
    
    // Connect to WebSocket
    err := wsClient.Connect()
    if err != nil {
        log.Fatal(err)
    }
    defer wsClient.Disconnect()
    
    // Subscribe to BTCUSDT ticker
    wsClient.SubscribeTicker("BTCUSDT", func(ticker *spot.Ticker24hr) {
        fmt.Printf("BTCUSDT: %s (Change: %s%%)\n", 
            ticker.LastPrice.String(), ticker.PriceChangePercent.String())
    })
    
    // Subscribe to BTCUSDT trades
    wsClient.SubscribeTrade("BTCUSDT", func(trade *spot.Trade) {
        side := "SELL"
        if trade.IsBuyerMaker {
            side = "BUY"
        }
        fmt.Printf("Trade: %s %s @ %s\n", side, trade.Qty.String(), trade.Price.String())
    })
    
    // Subscribe to BTCUSDT klines
    wsClient.SubscribeKline("BTCUSDT", spot.Interval1m, func(kline *spot.Kline) {
        fmt.Printf("Kline: O=%s H=%s L=%s C=%s\n",
            kline.Open.String(), kline.High.String(), 
            kline.Low.String(), kline.Close.String())
    })
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
}
```

## API Reference

### Spot Trading

#### Market Data
- `Ping()` - Test connectivity
- `GetServerTime()` - Get server time
- `GetExchangeInfo()` - Get exchange information
- `GetOrderBook(symbol, limit)` - Get order book
- `GetRecentTrades(symbol, limit)` - Get recent trades
- `GetKlines(symbol, interval, startTime, endTime, limit)` - Get kline data
- `GetTicker24hr(symbol)` - Get 24hr ticker
- `GetPrice(symbol)` - Get latest price
- `GetBookTicker(symbol)` - Get best bid/ask

#### Trading
- `NewOrder(req)` - Place new order
- `CancelOrder(symbol, orderID, origClientOrderID)` - Cancel order
- `GetOrder(symbol, orderID, origClientOrderID)` - Get order
- `GetOpenOrders(symbol)` - Get open orders
- `GetAllOrders(symbol, orderID, startTime, endTime, limit)` - Get all orders

#### Account
- `GetAccount()` - Get account information
- `GetUserTrades(symbol, orderID, startTime, endTime, fromID, limit)` - Get user trades
- `Transfer(req)` - Transfer between spot and futures
- `GetCommissionRate(symbol)` - Get commission rates

#### WebSocket Streams
- `NewWebSocketClient(testnet)` - Create WebSocket client
- `SubscribeTicker(symbol, handler)` - Subscribe to ticker stream
- `SubscribeMiniTicker(symbol, handler)` - Subscribe to mini ticker stream
- `SubscribeBookTicker(symbol, handler)` - Subscribe to book ticker stream
- `SubscribeTrade(symbol, handler)` - Subscribe to trade stream
- `SubscribeAggTrade(symbol, handler)` - Subscribe to aggregated trade stream
- `SubscribeKline(symbol, interval, handler)` - Subscribe to kline stream
- `SubscribeDepth(symbol, levels, handler)` - Subscribe to depth stream
- `SubscribeAllTickers(handler)` - Subscribe to all tickers stream

### Futures Trading

#### Market Data
- `Ping()` - Test connectivity
- `GetServerTime()` - Get server time
- `GetExchangeInfo()` - Get exchange information
- `GetOrderBook(symbol, limit)` - Get order book
- `GetRecentTrades(symbol, limit)` - Get recent trades
- `GetKlines(symbol, interval, startTime, endTime, limit)` - Get kline data
- `GetMarkPrice(symbol)` - Get mark price
- `GetTicker24hr(symbol)` - Get 24hr ticker
- `GetFundingRateHistory(symbol, startTime, endTime, limit)` - Get funding rate history

#### Trading
- `NewOrder(req)` - Place new order
- `PlaceMultipleOrders(orders)` - Place multiple orders
- `CancelOrder(symbol, orderID, origClientOrderID)` - Cancel order
- `CancelAllOpenOrders(symbol)` - Cancel all open orders
- `GetOrder(symbol, orderID, origClientOrderID)` - Get order
- `GetOpenOrders(symbol)` - Get open orders
- `GetAllOrders(symbol, orderID, startTime, endTime, limit)` - Get all orders

#### Account & Positions
- `GetAccount()` - Get account information
- `GetBalance()` - Get futures account balance
- `GetPositionInfo(symbol)` - Get position information
- `GetUserTrades(symbol, orderID, startTime, endTime, fromID, limit)` - Get user trades
- `Transfer(req)` - Transfer between spot and futures

#### Position Management
- `ChangePositionMode(req)` - Change position mode
- `GetCurrentPositionMode()` - Get current position mode
- `ChangeMultiAssetsMode(req)` - Change multi-assets mode
- `ChangeLeverage(req)` - Change leverage
- `ChangeMarginType(req)` - Change margin type

#### WebSocket Streams
- `NewWebSocketClient(testnet)` - Create WebSocket client
- `SubscribeTicker(symbol, handler)` - Subscribe to ticker stream
- `SubscribeMiniTicker(symbol, handler)` - Subscribe to mini ticker stream
- `SubscribeBookTicker(symbol, handler)` - Subscribe to book ticker stream
- `SubscribeTrade(symbol, handler)` - Subscribe to trade stream
- `SubscribeAggTrade(symbol, handler)` - Subscribe to aggregated trade stream
- `SubscribeKline(symbol, interval, handler)` - Subscribe to kline stream
- `SubscribeDepth(symbol, levels, handler)` - Subscribe to depth stream
- `SubscribeMarkPrice(symbol, handler)` - Subscribe to mark price stream
- `SubscribeAllMarkPrices(handler)` - Subscribe to all mark prices stream
- `SubscribeFundingRate(symbol, handler)` - Subscribe to funding rate stream
- `SubscribeAllTickers(handler)` - Subscribe to all tickers stream

## Configuration

### Client Configuration

```go
import "github.com/yiplee/aster-go/common"

config := &common.ClientConfig{
    APIKey:     "your-api-key",
    SecretKey:  "your-secret-key",
    BaseURL:    "https://sapi.asterdex.com", // or https://fapi.asterdex.com for futures
    Testnet:    false,
    Timeout:    30 * time.Second,
    RecvWindow: 5000,
}

// Create clients
spotClient := spot.NewClient(config)
futuresClient := futures.NewClient(config)
```

### Testnet Support

```go
// Spot testnet
spotClient := spot.NewClient(nil)
spotClient.SetTestnet(true)

// Futures testnet
futuresClient := futures.NewClient(nil)
futuresClient.SetTestnet(true)
```

## Decimal Precision

All price and quantity values use `decimal.Decimal` for high precision arithmetic:

```go
import "github.com/shopspring/decimal"

// Create decimal values
price := decimal.NewFromString("50000.123456789")
quantity := decimal.NewFromFloat(0.001)

// Perform calculations
total := price.Mul(quantity)

// Convert to string
priceStr := price.String() // "50000.123456789"
```

## Error Handling

The SDK provides comprehensive error handling:

```go
order, err := client.NewOrder(req)
if err != nil {
    // Check for API errors
    if apiErr, ok := err.(common.APIError); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.Code, apiErr.Msg)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Testing

The SDK includes comprehensive tests with mock implementations:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./spot
go test ./futures
go test ./common
```

## Examples

See the `examples/` directory for complete usage examples:

- `examples/spot_trading.go` - Spot trading examples
- `examples/futures_trading.go` - Futures trading examples
- `examples/market_data.go` - Market data examples

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue on GitHub
- Check the API documentation at [Aster API Docs](https://github.com/asterdex/api-docs)

## Changelog

### v1.0.0
- Initial release
- Spot trading support
- Futures trading support
- Decimal precision for all calculations
- Comprehensive test coverage
- Full API documentation
