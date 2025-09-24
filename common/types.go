package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// APIError represents an API error response
type APIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Msg)
}

// ClientConfig holds the configuration for the API client
type ClientConfig struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	Testnet    bool
	Timeout    time.Duration
	RecvWindow int64
	httpClient HTTPClient
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		BaseURL:    "https://sapi.asterdex.com",
		Testnet:    false,
		Timeout:    30 * time.Second,
		RecvWindow: 5000,
	}
}

// SignRequest creates a signature for authenticated requests
func SignRequest(params map[string]string, secretKey string) string {
	// Create query string
	var queryParts []string
	for key, value := range params {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
	}

	// Sort parameters
	sort.Strings(queryParts)
	queryString := strings.Join(queryParts, "&")

	// Create HMAC SHA256 signature
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

// BuildQueryString builds a query string from parameters
func BuildQueryString(params map[string]any) string {
	values := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case string:
			if v != "" {
				values.Add(key, v)
			}
		case int:
			if v != 0 {
				values.Add(key, strconv.Itoa(v))
			}
		case int64:
			if v != 0 {
				values.Add(key, strconv.FormatInt(v, 10))
			}
		case float64:
			if v != 0 {
				values.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
			}
		case bool:
			values.Add(key, strconv.FormatBool(v))
		}
	}
	return values.Encode()
}

// GetTimestamp returns the current timestamp in milliseconds
func GetTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FormatFloat formats a float64 to string with specified precision
func FormatFloat(f float64, precision int) string {
	return strconv.FormatFloat(f, 'f', precision, 64)
}

// ParseFloat parses a string to float64
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseInt parses a string to int64
func ParseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
