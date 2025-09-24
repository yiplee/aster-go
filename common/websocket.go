package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient represents a WebSocket client
type WebSocketClient struct {
	conn              *websocket.Conn
	url               string
	mu                sync.RWMutex
	handlers          map[string][]func(json.RawMessage)
	connected         bool
	reconnect         bool
	reconnectInterval time.Duration
	ctx               context.Context
	cancel            context.CancelFunc
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(baseURL string) *WebSocketClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketClient{
		url:               baseURL,
		handlers:          make(map[string][]func(json.RawMessage)),
		reconnect:         true,
		reconnectInterval: 5 * time.Second,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Connect establishes a WebSocket connection
func (c *WebSocketClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	u, err := url.Parse(c.url)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.conn = conn
	c.connected = true

	// Start message handler
	go c.handleMessages()

	return nil
}

// Disconnect closes the WebSocket connection
func (c *WebSocketClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()
	c.connected = false

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// Subscribe subscribes to a stream
func (c *WebSocketClient) Subscribe(stream string, handler func(json.RawMessage)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.handlers[stream] == nil {
		c.handlers[stream] = make([]func(json.RawMessage), 0)
	}
	c.handlers[stream] = append(c.handlers[stream], handler)

	// Send subscription message if connected
	if c.connected && c.conn != nil {
		subscribeMsg := map[string]any{
			"method": "SUBSCRIBE",
			"params": []string{stream},
			"id":     time.Now().Unix(),
		}

		c.conn.WriteJSON(subscribeMsg)
	}
}

// Unsubscribe unsubscribes from a stream
func (c *WebSocketClient) Unsubscribe(stream string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.handlers, stream)

	// Send unsubscription message if connected
	if c.connected && c.conn != nil {
		unsubscribeMsg := map[string]any{
			"method": "UNSUBSCRIBE",
			"params": []string{stream},
			"id":     time.Now().Unix(),
		}

		c.conn.WriteJSON(unsubscribeMsg)
	}
}

// handleMessages handles incoming WebSocket messages
func (c *WebSocketClient) handleMessages() {
	defer func() {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()

		// Attempt to reconnect if enabled
		if c.reconnect {
			time.Sleep(c.reconnectInterval)
			c.Connect()
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Printf("WebSocket read error: %v\n", err)
				return
			}

			var wsMsg WebSocketMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				// Try to parse as direct data
				c.handleRawMessage(message)
				continue
			}

			c.handleStreamMessage(wsMsg.Stream, wsMsg.Data)
		}
	}
}

// handleRawMessage handles raw messages
func (c *WebSocketClient) handleRawMessage(message []byte) {
	// Try to parse as WebSocketMessage
	var wsMsg WebSocketMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		return
	}

	// Handle the stream message
	c.handleStreamMessage(wsMsg.Stream, wsMsg.Data)
}

// handleStreamMessage handles messages for a specific stream
func (c *WebSocketClient) handleStreamMessage(stream string, data json.RawMessage) {
	c.mu.RLock()
	handlers := c.handlers[stream]
	c.mu.RUnlock()

	for _, handler := range handlers {
		go handler(data)
	}
}

// SetReconnect sets the reconnect behavior
func (c *WebSocketClient) SetReconnect(reconnect bool, interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reconnect = reconnect
	c.reconnectInterval = interval
}

// IsConnected returns whether the client is connected
func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}
