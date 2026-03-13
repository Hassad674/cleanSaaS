package ws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeConn is not needed — we test the Hub's channel-based logic directly
// by creating clients with nil conns and intercepting their send channels.

func newTestHub(t *testing.T) *Hub {
	t.Helper()
	h := NewHub()
	go h.Run()
	t.Cleanup(func() { h.Stop() })
	return h
}

// testClient creates a Client with a send channel but no real WebSocket conn.
// This is safe because hub tests only exercise the channel-based routing logic,
// not the actual WebSocket read/write pumps.
func testClient(hub *Hub, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   nil,
		userID: userID,
		send:   make(chan []byte, sendBufferSize),
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	h := newTestHub(t)

	client := testClient(h, "user-1")
	h.register <- client

	// Give the hub goroutine time to process.
	time.Sleep(20 * time.Millisecond)

	assert.Equal(t, 1, h.ConnectedUserCount())

	h.unregister <- client

	time.Sleep(20 * time.Millisecond)

	assert.Equal(t, 0, h.ConnectedUserCount())
}

func TestHub_MultipleConnectionsSameUser(t *testing.T) {
	h := newTestHub(t)

	c1 := testClient(h, "user-1")
	c2 := testClient(h, "user-1")

	h.register <- c1
	h.register <- c2

	time.Sleep(20 * time.Millisecond)

	// Still one unique user.
	assert.Equal(t, 1, h.ConnectedUserCount())

	// Unregister one — user still has a connection.
	h.unregister <- c1
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, h.ConnectedUserCount())

	// Unregister the last — user fully disconnected.
	h.unregister <- c2
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 0, h.ConnectedUserCount())
}

func TestHub_SendToUser_DeliversToCorrectClient(t *testing.T) {
	h := newTestHub(t)

	c1 := testClient(h, "user-1")
	c2 := testClient(h, "user-2")

	h.register <- c1
	h.register <- c2
	time.Sleep(20 * time.Millisecond)

	msg := []byte(`{"type":"notification","payload":{"title":"hello"}}`)
	err := h.SendToUser("user-1", msg)
	require.NoError(t, err)

	// user-1 should receive the message.
	select {
	case received := <-c1.send:
		assert.Equal(t, msg, received)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("user-1 did not receive message")
	}

	// user-2 should NOT receive the message.
	select {
	case <-c2.send:
		t.Fatal("user-2 should not have received message")
	case <-time.After(50 * time.Millisecond):
		// expected
	}
}

func TestHub_SendToUser_MultipleConnections(t *testing.T) {
	h := newTestHub(t)

	c1 := testClient(h, "user-1")
	c2 := testClient(h, "user-1")

	h.register <- c1
	h.register <- c2
	time.Sleep(20 * time.Millisecond)

	msg := []byte(`{"type":"system","payload":{}}`)
	err := h.SendToUser("user-1", msg)
	require.NoError(t, err)

	// Both connections should receive the message.
	for _, c := range []*Client{c1, c2} {
		select {
		case received := <-c.send:
			assert.Equal(t, msg, received)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("client did not receive message")
		}
	}
}

func TestHub_SendToUser_NoConnectionsIsNoOp(t *testing.T) {
	h := newTestHub(t)

	err := h.SendToUser("nonexistent-user", []byte(`{}`))
	assert.NoError(t, err)
}

func TestHub_Broadcast(t *testing.T) {
	h := newTestHub(t)

	c1 := testClient(h, "user-1")
	c2 := testClient(h, "user-2")
	c3 := testClient(h, "user-3")

	h.register <- c1
	h.register <- c2
	h.register <- c3
	time.Sleep(20 * time.Millisecond)

	msg := []byte(`{"type":"system","payload":{"message":"server restart"}}`)
	h.Broadcast(msg)

	for _, c := range []*Client{c1, c2, c3} {
		select {
		case received := <-c.send:
			assert.Equal(t, msg, received)
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("client %s did not receive broadcast", c.userID)
		}
	}
}

func TestHub_BroadcastMessage(t *testing.T) {
	h := newTestHub(t)

	c := testClient(h, "user-1")
	h.register <- c
	time.Sleep(20 * time.Millisecond)

	m, err := NewMessage("system", map[string]string{"event": "maintenance"})
	require.NoError(t, err)

	err = h.BroadcastMessage(m)
	require.NoError(t, err)

	select {
	case received := <-c.send:
		var msg Message
		err := json.Unmarshal(received, &msg)
		require.NoError(t, err)
		assert.Equal(t, "system", msg.Type)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("did not receive broadcast message")
	}
}

func TestHub_Stop_ClosesAllClients(t *testing.T) {
	h := NewHub()
	go h.Run()

	c1 := testClient(h, "user-1")
	c2 := testClient(h, "user-2")

	h.register <- c1
	h.register <- c2
	time.Sleep(20 * time.Millisecond)

	h.Stop()

	// After stop, send channels should be closed.
	_, ok1 := <-c1.send
	assert.False(t, ok1, "client 1 send channel should be closed")

	_, ok2 := <-c2.send
	assert.False(t, ok2, "client 2 send channel should be closed")

	assert.Equal(t, 0, h.ConnectedUserCount())
}

func TestNewMessage(t *testing.T) {
	payload := map[string]string{"key": "value"}
	msg, err := NewMessage("test", payload)
	require.NoError(t, err)
	assert.Equal(t, "test", msg.Type)

	var decoded map[string]string
	err = json.Unmarshal(msg.Payload, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "value", decoded["key"])
}

func TestNewMessage_InvalidPayload(t *testing.T) {
	// Functions cannot be marshalled to JSON.
	_, err := NewMessage("test", func() {})
	assert.Error(t, err)
}
