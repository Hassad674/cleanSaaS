package ws

import "encoding/json"

// Message represents a WebSocket message exchanged between client and server.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewMessage creates a new Message by marshalling the given payload to JSON.
func NewMessage(msgType string, payload any) (*Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:    msgType,
		Payload: data,
	}, nil
}
