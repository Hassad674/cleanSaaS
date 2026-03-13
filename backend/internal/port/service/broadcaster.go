package service

// Broadcaster defines the interface for sending real-time messages to connected users.
// This is implemented by the WebSocket hub, but the interface lives in port/
// to maintain hexagonal architecture — app services depend on this abstraction,
// not on the concrete WebSocket implementation.
type Broadcaster interface {
	// SendToUser sends a raw JSON message to all active connections of the given user.
	// Returns nil if the user has no active connections (fire-and-forget semantics).
	SendToUser(userID string, msg []byte) error
}
