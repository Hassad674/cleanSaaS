package ws

import (
	"encoding/json"
	"log/slog"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	// Registered clients, keyed by userID for efficient targeted delivery.
	clients map[string]map[*Client]struct{}

	// Register requests from clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Broadcast sends a message to all connected clients.
	broadcast chan []byte

	// stop signals the Run goroutine to shut down.
	stop chan struct{}

	// done is closed when Run has fully stopped.
	done chan struct{}

	mu sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
	}
}

// Run starts the hub's main event loop. It should be called in its own goroutine.
func (h *Hub) Run() {
	defer close(h.done)

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]struct{})
			}
			h.clients[client.userID][client] = struct{}{}
			h.mu.Unlock()

			slog.Info("ws: client connected",
				slog.String("user_id", client.userID),
				slog.Int("user_connections", len(h.clients[client.userID])),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.userID]; ok {
				if _, exists := conns[client]; exists {
					delete(conns, client)
					close(client.send)
					if len(conns) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()

			slog.Info("ws: client disconnected", slog.String("user_id", client.userID))

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, conns := range h.clients {
				for client := range conns {
					select {
					case client.send <- message:
					default:
						// Client send buffer full — drop connection.
						go h.removeClient(client)
					}
				}
			}
			h.mu.RUnlock()

		case <-h.stop:
			h.mu.Lock()
			for _, conns := range h.clients {
				for client := range conns {
					close(client.send)
				}
			}
			h.clients = make(map[string]map[*Client]struct{})
			h.mu.Unlock()
			return
		}
	}
}

// Stop gracefully shuts down the hub, closing all client connections.
func (h *Hub) Stop() {
	close(h.stop)
	<-h.done
}

// Register adds a client to the hub.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendToUser sends a message to all connections of a specific user.
// Returns nil even if the user has no active connections (fire-and-forget).
func (h *Hub) SendToUser(userID string, msg []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, ok := h.clients[userID]
	if !ok {
		return nil
	}

	for client := range conns {
		select {
		case client.send <- msg:
		default:
			go h.removeClient(client)
		}
	}
	return nil
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// BroadcastMessage marshals a Message and sends it to all connected clients.
func (h *Hub) BroadcastMessage(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	h.Broadcast(data)
	return nil
}

// ConnectedUserCount returns the number of unique users with active connections.
func (h *Hub) ConnectedUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// removeClient unregisters a client via the channel (safe for concurrent use).
func (h *Hub) removeClient(client *Client) {
	h.unregister <- client
}
