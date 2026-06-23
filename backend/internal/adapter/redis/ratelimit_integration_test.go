//go:build integration

// Integration tests for the Redis-backed rate limiter and pub/sub broadcaster.
// They run against a LIVE Redis and degrade gracefully: if REDIS_URL is unset or
// Redis is unreachable, each test t.Skip()s rather than failing, so a CI lane
// without Redis still passes.
//
//	docker compose up -d redis
//	REDIS_URL=redis://localhost:6379 go test -tags=integration ./internal/adapter/redis/... -count=1
package redis

import (
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRedisURL() string {
	if v := os.Getenv("REDIS_URL"); v != "" {
		return v
	}
	return "redis://localhost:6379"
}

func dialOrSkip(t *testing.T) *goredis.Client {
	t.Helper()
	client, err := Connect(testRedisURL())
	if err != nil {
		t.Skipf("skipping redis integration test: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestRedisRateLimiter_EnforcesLimitWithinWindow proves the limiter rejects once
// the per-window budget is spent for a key.
func TestRedisRateLimiter_EnforcesLimitWithinWindow(t *testing.T) {
	client := dialOrSkip(t)
	key := "itest-key-" + time.Now().Format("150405.000000")

	rl := NewRateLimiter(client, 5, "itest-enforce", quietLogger())

	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow(key), "request %d within budget should be allowed", i)
	}
	assert.False(t, rl.Allow(key), "6th request over the budget must be rejected")
}

// TestRedisRateLimiter_SharedAcrossInstances is the headline guarantee: two
// limiters built from SEPARATE clients (two app instances) but the SAME keyspace
// and key share one budget in Redis. Combined they must allow exactly `limit`
// requests, not `limit` per instance — which is what fixes the in-memory
// "limits multiply per instance" bug.
func TestRedisRateLimiter_SharedAcrossInstances(t *testing.T) {
	clientA := dialOrSkip(t)
	clientB := dialOrSkip(t)

	keyspace := "itest-shared"
	key := "shared-key-" + time.Now().Format("150405.000000")
	const limit = 10

	a := NewRateLimiter(clientA, limit, keyspace, quietLogger())
	b := NewRateLimiter(clientB, limit, keyspace, quietLogger())

	allowed := 0
	// Alternate between the two "instances"; the shared Redis counter must cap the
	// total at `limit` regardless of which instance serves a request.
	for i := 0; i < limit*2; i++ {
		var ok bool
		if i%2 == 0 {
			ok = a.Allow(key)
		} else {
			ok = b.Allow(key)
		}
		if ok {
			allowed++
		}
	}

	assert.Equal(t, limit, allowed,
		"a shared Redis budget must allow exactly %d total across both instances", limit)
}

// TestRedisBroadcaster_CrossInstanceDelivery proves WebSocket pub/sub fan-out:
// a message published via instance A's broadcaster is delivered to instance B's
// LOCAL hub for the target user.
func TestRedisBroadcaster_CrossInstanceDelivery(t *testing.T) {
	clientA := dialOrSkip(t)
	clientB := dialOrSkip(t)

	// recordingHub captures local deliveries on "instance B".
	hubB := &recordingHub{got: make(chan delivered, 4)}

	// Broadcaster A publishes; broadcaster B subscribes + delivers to hubB.
	bcastA := NewBroadcaster(clientA, &recordingHub{got: make(chan delivered, 4)}, quietLogger())
	t.Cleanup(bcastA.Stop)
	bcastB := NewBroadcaster(clientB, hubB, quietLogger())
	t.Cleanup(bcastB.Stop)

	// Give B's subscriber a moment to establish the subscription.
	time.Sleep(100 * time.Millisecond)

	payload := []byte(`{"type":"notification","payload":{"title":"hi"}}`)
	require.NoError(t, bcastA.SendToUser("user-42", payload))

	select {
	case d := <-hubB.got:
		assert.Equal(t, "user-42", d.userID)
		assert.Equal(t, payload, d.msg)
	case <-time.After(2 * time.Second):
		t.Fatal("instance B did not receive the cross-instance message")
	}
}

// delivered records one local SendToUser call.
type delivered struct {
	userID string
	msg    []byte
}

// recordingHub is a LocalDelivery test double that records what it was asked to
// deliver locally, standing in for the real WebSocket hub.
type recordingHub struct {
	mu  sync.Mutex
	got chan delivered
}

func (h *recordingHub) SendToUser(userID string, msg []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	cp := make([]byte, len(msg))
	copy(cp, msg)
	h.got <- delivered{userID: userID, msg: cp}
	return nil
}
