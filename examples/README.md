# Aster SDK Examples

This directory contains example applications demonstrating how to use the Aster SDK for various trading operations.

## Available Examples

### 1. Spot Trading (`spot_trading/`)
Demonstrates spot market operations including:
- Server connectivity testing
- Market data retrieval
- Order management
- Account information

**Run:** `cd spot_trading && go run main.go`

### 2. Futures Trading (`futures_trading/`)
Demonstrates futures market operations including:
- Futures market data
- Position management
- Funding rate information
- Leverage management

**Run:** `cd futures_trading && go run main.go`

### 3. Market Data (`market_data/`)
Comprehensive market data retrieval for both spot and futures:
- Real-time tickers and prices
- Order book data
- Historical data
- Kline intervals

**Run:** `cd market_data && go run main.go`

### 4. WebSocket Example (`websocket_example/`)
Real-time data streaming using WebSocket connections:
- Live ticker updates
- Order book streams
- Trade streams
- Mark price updates

**Run:** `cd websocket_example && go run main.go`

## Prerequisites

- Go 1.24.2 or later
- Valid API credentials (for authenticated operations)

## Configuration

Before running examples that require authentication, set your API credentials:

```go
client.SetAPIKey("your-api-key", "your-secret-key")
```

## Safety Note

Order placement examples are commented out by default to prevent accidental trades. Uncomment and modify them only when you're ready to place real orders.
