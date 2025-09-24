package spot

import (
	"testing"
)

func TestSymbolValidation(t *testing.T) {
	client := NewClient(nil)

	// Test GetTicker24hr with empty symbol
	_, err := client.GetTicker24hr("")
	if err == nil {
		t.Error("Expected error for empty symbol in GetTicker24hr")
	}
	if err.Error() != "symbol is required" {
		t.Errorf("Expected 'symbol is required', got '%s'", err.Error())
	}

	// Test GetPrice with empty symbol
	_, err = client.GetPrice("")
	if err == nil {
		t.Error("Expected error for empty symbol in GetPrice")
	}
	if err.Error() != "symbol is required" {
		t.Errorf("Expected 'symbol is required', got '%s'", err.Error())
	}

	// Test GetBookTicker with empty symbol
	_, err = client.GetBookTicker("")
	if err == nil {
		t.Error("Expected error for empty symbol in GetBookTicker")
	}
	if err.Error() != "symbol is required" {
		t.Errorf("Expected 'symbol is required', got '%s'", err.Error())
	}

	// Test GetOrderBook with empty symbol
	_, err = client.GetOrderBook("", 10)
	if err == nil {
		t.Error("Expected error for empty symbol in GetOrderBook")
	}
	if err.Error() != "symbol is required" {
		t.Errorf("Expected 'symbol is required', got '%s'", err.Error())
	}

	// Test GetRecentTrades with empty symbol
	_, err = client.GetRecentTrades("", 10)
	if err == nil {
		t.Error("Expected error for empty symbol in GetRecentTrades")
	}
	if err.Error() != "symbol is required" {
		t.Errorf("Expected 'symbol is required', got '%s'", err.Error())
	}
}
