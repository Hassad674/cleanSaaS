package handler

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ws"
)

// WSHandler handles WebSocket upgrade requests.
type WSHandler struct {
	hub      *ws.Hub
	jwtMaker *jwt.Maker
	upgrader websocket.Upgrader
}

// NewWSHandler creates a new WebSocket handler.
// allowedOrigins controls which origins can open WebSocket connections.
func NewWSHandler(hub *ws.Hub, jwtMaker *jwt.Maker, allowedOrigins ...string) *WSHandler {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = true
	}

	return &WSHandler{
		hub:      hub,
		jwtMaker: jwtMaker,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return originSet[origin]
			},
		},
	}
}

// Upgrade handles the GET /ws endpoint.
// Authentication is done via a query parameter ?token=xxx because the browser
// WebSocket API does not support custom Authorization headers.
func (h *WSHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token query parameter", http.StatusUnauthorized)
		return
	}

	claims, err := h.jwtMaker.Validate(token)
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws: upgrade failed", slog.String("error", err.Error()))
		return
	}

	client := ws.NewClient(h.hub, conn, claims.UserID)
	h.hub.Register(client)

	// Start read and write pumps in separate goroutines.
	go client.WritePump()
	go client.ReadPump()
}
