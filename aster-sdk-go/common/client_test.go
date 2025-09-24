package common

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of HTTPClient
type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func TestNewClient(t *testing.T) {
	config := DefaultConfig()
	client := NewClient(config)
	
	if client == nil {
		t.Error("Expected client to not be nil")
	}
	
	if client.GetConfig() != config {
		t.Error("Expected client config to match provided config")
	}
}

func TestNewClientWithNilConfig(t *testing.T) {
	client := NewClient(nil)
	
	if client == nil {
		t.Error("Expected client to not be nil")
	}
	
	config := client.GetConfig()
	if config.BaseURL != "https://sapi.asterdex.com" {
		t.Errorf("Expected default base URL, got %s", config.BaseURL)
	}
}

func TestSetAPIKey(t *testing.T) {
	client := NewClient(nil)
	apiKey := "test-api-key"
	secretKey := "test-secret-key"
	
	client.SetAPIKey(apiKey, secretKey)
	
	config := client.GetConfig()
	if config.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, config.APIKey)
	}
	if config.SecretKey != secretKey {
		t.Errorf("Expected secret key %s, got %s", secretKey, config.SecretKey)
	}
}

func TestSetBaseURL(t *testing.T) {
	client := NewClient(nil)
	newURL := "https://test.example.com"
	
	client.SetBaseURL(newURL)
	
	config := client.GetConfig()
	if config.BaseURL != newURL {
		t.Errorf("Expected base URL %s, got %s", newURL, config.BaseURL)
	}
}

func TestSetHTTPClient(t *testing.T) {
	client := NewClient(nil)
	mockClient := &MockHTTPClient{}
	
	client.SetHTTPClient(mockClient)
	
	// Test that the client was set by checking if we can call DoRequest
	// This is a basic test - in practice you'd want more comprehensive testing
	if client.httpClient != mockClient {
		t.Error("Expected HTTP client to be set to mock client")
	}
}

func TestDoRequest(t *testing.T) {
	client := NewClient(nil)
	client.SetBaseURL("https://api.example.com")
	
	// Create a mock response
	responseBody := `{"success": true}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}
	
	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)
	
	// Test GET request
	resp, err := client.DoRequest("GET", "/test", map[string]any{"param": "value"}, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestDoRequestWithSignature(t *testing.T) {
	client := NewClient(nil)
	client.SetAPIKey("test-api-key", "test-secret-key")
	client.SetBaseURL("https://api.example.com")
	
	// Create a mock response
	responseBody := `{"success": true}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}
	
	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)
	
	// Test signed request
	resp, err := client.DoRequest("POST", "/test", map[string]any{"param": "value"}, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestParseResponse(t *testing.T) {
	client := NewClient(nil)
	
	// Test successful response
	responseBody := `{"message": "success", "data": {"id": 123}}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}
	
	var result map[string]any
	err := client.ParseResponse(mockResponse, &result)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result["message"] != "success" {
		t.Errorf("Expected message 'success', got %v", result["message"])
	}
}

func TestParseResponseWithAPIError(t *testing.T) {
	client := NewClient(nil)
	
	// Test API error response
	responseBody := `{"code": -1121, "msg": "Invalid symbol."}`
	mockResponse := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}
	
	var result map[string]any
	err := client.ParseResponse(mockResponse, &result)
	if err == nil {
		t.Error("Expected error for API error response")
	}
	
	apiErr, ok := err.(APIError)
	if !ok {
		t.Errorf("Expected APIError, got %T", err)
	}
	
	if apiErr.Code != -1121 {
		t.Errorf("Expected error code -1121, got %d", apiErr.Code)
	}
}

func TestDo(t *testing.T) {
	client := NewClient(nil)
	client.SetBaseURL("https://api.example.com")
	
	// Create a mock response
	responseBody := `{"success": true}`
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
		Header:     make(http.Header),
	}
	
	mockClient := &MockHTTPClient{
		Response: mockResponse,
	}
	client.SetHTTPClient(mockClient)
	
	// Test Do method
	var result map[string]any
	err := client.Do("GET", "/test", map[string]any{"param": "value"}, &result, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result["success"] != true {
		t.Errorf("Expected success true, got %v", result["success"])
	}
}

func TestDoWithError(t *testing.T) {
	client := NewClient(nil)
	client.SetBaseURL("https://api.example.com")
	
	// Create a mock client that returns an error
	mockClient := &MockHTTPClient{
		Error: io.ErrUnexpectedEOF,
	}
	client.SetHTTPClient(mockClient)
	
	// Test Do method with error
	var result map[string]any
	err := client.Do("GET", "/test", map[string]any{"param": "value"}, &result, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestClientConfig(t *testing.T) {
	config := &ClientConfig{
		APIKey:     "test-key",
		SecretKey:  "test-secret",
		BaseURL:    "https://test.example.com",
		Testnet:    true,
		Timeout:    60 * time.Second,
		RecvWindow: 10000,
	}
	
	client := NewClient(config)
	
	if client.GetConfig() != config {
		t.Error("Expected client config to match provided config")
	}
}