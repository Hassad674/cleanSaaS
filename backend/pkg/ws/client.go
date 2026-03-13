package ws

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 30 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 4096

	// Send channel buffer size.
	sendBufferSize = 256
)

// Client represents a single WebSocket connection.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID string
	send   chan []byte
}

// NewClient creates a new client for the given WebSocket connection and user.
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, sendBufferSize),
	}
}

// UserID returns the authenticated user ID associated with this client.
func (c *Client) UserID() string {
	return c.userID
}

// ReadPump pumps messages from the WebSocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. It ensures
// that there is at most one reader on a connection by executing all reads
// from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error("ws: failed to set read deadline", slog.String("error", err.Error()))
		return
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				slog.Warn("ws: unexpected close", slog.String("error", err.Error()), slog.String("user_id", c.userID))
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			slog.Warn("ws: invalid message format", slog.String("error", err.Error()), slog.String("user_id", c.userID))
			continue
		}

		// For now, client-to-server messages are logged but not routed.
		// Features can extend this by adding message type handlers to the hub.
		slog.Debug("ws: received message", slog.String("type", msg.Type), slog.String("user_id", c.userID))
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
//
// A goroutine running WritePump is started for each connection. It ensures
// that there is at most one writer to a connection by executing all writes
// from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				slog.Error("ws: failed to set write deadline", slog.String("error", err.Error()))
				return
			}
			if !ok {
				// Hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				return
			}

			// Drain queued messages into the current write to reduce syscalls.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte("\n")); err != nil {
					break
				}
				if _, err := w.Write(<-c.send); err != nil {
					break
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				slog.Error("ws: failed to set write deadline for ping", slog.String("error", err.Error()))
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
