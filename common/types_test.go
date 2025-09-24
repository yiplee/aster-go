package common

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.BaseURL != "https://sapi.asterdex.com" {
		t.Errorf("Expected base URL to be 'https://sapi.asterdex.com', got '%s'", config.BaseURL)
	}
	
	if config.Testnet {
		t.Error("Expected testnet to be false")
	}
	
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", config.Timeout)
	}
	
	if config.RecvWindow != 5000 {
		t.Errorf("Expected recvWindow to be 5000, got %d", config.RecvWindow)
	}
}

func TestSignRequest(t *testing.T) {
	params := map[string]string{
		"symbol":    "BTCUSDT",
		"side":      "BUY",
		"type":      "LIMIT",
		"quantity":  "1.0",
		"price":     "50000",
		"timestamp": "1640995200000",
	}
	
	secretKey := "test-secret-key"
	signature := SignRequest(params, secretKey)
	
	if signature == "" {
		t.Error("Expected signature to not be empty")
	}
	
	// Test that the same parameters produce the same signature
	signature2 := SignRequest(params, secretKey)
	if signature != signature2 {
		t.Error("Expected same parameters to produce same signature")
	}
	
	// Test that different parameters produce different signatures
	params["price"] = "51000"
	signature3 := SignRequest(params, secretKey)
	if signature == signature3 {
		t.Error("Expected different parameters to produce different signature")
	}
}

func TestBuildQueryString(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]any
		expected string
	}{
		{
			name:     "empty params",
			params:   map[string]any{},
			expected: "",
		},
		{
			name: "string params",
			params: map[string]any{
				"symbol": "BTCUSDT",
				"side":   "BUY",
			},
			expected: "side=BUY&symbol=BTCUSDT",
		},
		{
			name: "mixed types",
			params: map[string]any{
				"symbol":   "BTCUSDT",
				"limit":    100,
				"active":   true,
				"price":    50000.5,
			},
			expected: "active=true&limit=100&price=50000.5&symbol=BTCUSDT",
		},
		{
			name: "zero values",
			params: map[string]any{
				"symbol": "BTCUSDT",
				"limit":  0,
				"price":  0.0,
			},
			expected: "symbol=BTCUSDT",
		},
		{
			name: "empty string",
			params: map[string]any{
				"symbol": "BTCUSDT",
				"side":   "",
			},
			expected: "symbol=BTCUSDT",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildQueryString(tt.params)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetTimestamp(t *testing.T) {
	before := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := GetTimestamp()
	after := time.Now().UnixNano() / int64(time.Millisecond)
	
	if timestamp < before || timestamp > after {
		t.Errorf("Timestamp %d should be between %d and %d", timestamp, before, after)
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		value    float64
		precision int
		expected string
	}{
		{1.23456789, 2, "1.23"},
		{1.23456789, 4, "1.2346"},
		{0.0, 2, "0.00"},
		{123.456, 0, "123"},
		{123.456, 1, "123.5"},
	}
	
	for _, tt := range tests {
		result := FormatFloat(tt.value, tt.precision)
		if result != tt.expected {
			t.Errorf("FormatFloat(%f, %d) = %s, expected %s", tt.value, tt.precision, result, tt.expected)
		}
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"1.23", 1.23, false},
		{"0", 0.0, false},
		{"123.456", 123.456, false},
		{"invalid", 0.0, true},
		{"", 0.0, true},
	}
	
	for _, tt := range tests {
		result, err := ParseFloat(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseFloat(%s) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseFloat(%s) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseFloat(%s) = %f, expected %f", tt.input, result, tt.expected)
			}
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"-456", -456, false},
		{"invalid", 0, true},
		{"", 0, true},
	}
	
	for _, tt := range tests {
		result, err := ParseInt(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseInt(%s) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseInt(%s) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseInt(%s) = %d, expected %d", tt.input, result, tt.expected)
			}
		}
	}
}