package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents the base API client
type Client struct {
	config     *ClientConfig
	httpClient HTTPClient
}

// NewClient creates a new API client
func NewClient(config *ClientConfig) *Client {
	if config == nil {
		config = DefaultConfig()
	}
	
	if config.httpClient == nil {
		config.httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}
	
	return &Client{
		config:     config,
		httpClient: config.httpClient,
	}
}

// SetHTTPClient sets a custom HTTP client
func (c *Client) SetHTTPClient(client HTTPClient) {
	c.httpClient = client
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *ClientConfig {
	return c.config
}

// SetAPIKey sets the API key and secret
func (c *Client) SetAPIKey(apiKey, secretKey string) {
	c.config.APIKey = apiKey
	c.config.SecretKey = secretKey
}

// SetBaseURL sets the base URL
func (c *Client) SetBaseURL(baseURL string) {
	c.config.BaseURL = baseURL
}

// DoRequest performs an HTTP request
func (c *Client) DoRequest(method, endpoint string, params map[string]any, signed bool) (*http.Response, error) {
	requestURL := c.config.BaseURL + endpoint
	
	var body io.Reader
	var queryString string
	
	if method == "GET" || method == "DELETE" {
		// Add parameters to query string
		if len(params) > 0 {
			queryString = BuildQueryString(params)
			requestURL += "?" + queryString
		}
	} else {
		// Add parameters to request body
		if len(params) > 0 {
			formData := BuildQueryString(params)
			body = bytes.NewBufferString(formData)
		}
	}
	
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}
	
	// Set headers
	if c.config.APIKey != "" {
		req.Header.Set("X-MBX-APIKEY", c.config.APIKey)
	}
	
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	
	// Add signature for signed requests
	if signed && c.config.SecretKey != "" {
		signParams := make(map[string]string)
		for k, v := range params {
			signParams[k] = fmt.Sprintf("%v", v)
		}
		
		// Add timestamp
		timestamp := GetTimestamp()
		signParams["timestamp"] = fmt.Sprintf("%d", timestamp)
		
		// Add recvWindow if not present
		if _, exists := signParams["recvWindow"]; !exists {
			signParams["recvWindow"] = fmt.Sprintf("%d", c.config.RecvWindow)
		}
		
		signature := SignRequest(signParams, c.config.SecretKey)
		
		if method == "GET" || method == "DELETE" {
			requestURL += "&signature=" + signature
			req.URL, _ = url.Parse(requestURL)
		} else {
			// Add signature to form data
			formData := BuildQueryString(params)
			formData += "&timestamp=" + fmt.Sprintf("%d", timestamp)
			formData += "&recvWindow=" + fmt.Sprintf("%d", c.config.RecvWindow)
			formData += "&signature=" + signature
			body = bytes.NewBufferString(formData)
			req.Body = io.NopCloser(body)
		}
	}
	
	return c.httpClient.Do(req)
}

// ParseResponse parses the HTTP response
func (c *Client) ParseResponse(resp *http.Response, result any) error {
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	// Check for API errors
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return apiErr
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse successful response
	if result != nil {
		return json.Unmarshal(body, result)
	}
	
	return nil
}

// Do performs a request and parses the response
func (c *Client) Do(method, endpoint string, params map[string]any, result any, signed bool) error {
	resp, err := c.DoRequest(method, endpoint, params, signed)
	if err != nil {
		return err
	}
	
	return c.ParseResponse(resp, result)
}