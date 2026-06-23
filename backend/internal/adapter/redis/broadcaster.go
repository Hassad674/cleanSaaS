package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// broadcastChannel is the Redis pub/sub channel every instance publishes to and
// subscribes on for cross-instance WebSocket fan-out.
const broadcastChannel = "ws:send-to-user"

// publishTimeout bounds a single publish so a slow/unreachable Redis can never
// block the notification path; on failure the publish is logged, never blocked on.
const publishTimeout = 200 * time.Millisecond

// LocalDelivery is the slice of the WebSocket hub the broadcaster needs: deliver a
// message to the connections of userID that live on THIS instance. The hub's
// existing SendToUser satisfies it, so its concurrency model is untouched.
type LocalDelivery interface {
	SendToUser(userID string, msg []byte) error
}

// userMessage is the wire envelope published to the broadcast channel so every
// subscribing instance knows which user a payload is for.
type userMessage struct {
	UserID  string `json:"user_id"`
	Payload []byte `json:"payload"`
}

// Broadcaster fans WebSocket messages out across instances using Redis pub/sub.
// SendToUser publishes to a shared channel instead of delivering only to local
// sockets; a background subscriber on EVERY instance receives every published
// message and delivers it to its own local connections for that user. This makes
// SendToUser reach a user's sockets no matter which instance they are connected
// to. It implements service.Broadcaster (same SendToUser signature as the hub),
// so it is a drop-in replacement at the wiring site.
type Broadcaster struct {
	client *goredis.Client
	local  LocalDelivery
	pubsub *goredis.PubSub
	logger *slog.Logger
	stop   chan struct{}
	done   chan struct{}
}

// NewBroadcaster wraps the local hub with Redis pub/sub fan-out and starts the
// subscriber goroutine. local is the in-process hub (its SendToUser does the
// actual socket writes on this instance).
func NewBroadcaster(client *goredis.Client, local LocalDelivery, logger *slog.Logger) *Broadcaster {
	if logger == nil {
		logger = slog.Default()
	}
	b := &Broadcaster{
		client: client,
		local:  local,
		pubsub: client.Subscribe(context.Background(), broadcastChannel),
		logger: logger,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	go b.subscribe()
	return b
}

// SendToUser publishes the message to the shared channel. Delivery to the actual
// sockets (on this and every other instance) happens in subscribe(), so the
// originating instance does NOT also deliver directly — that would double-send to
// locally-connected sockets. Fire-and-forget: a publish error is logged, never
// returned to the caller's critical path.
func (b *Broadcaster) SendToUser(userID string, msg []byte) error {
	envelope, err := json.Marshal(userMessage{UserID: userID, Payload: msg})
	if err != nil {
		return fmt.Errorf("marshalling ws envelope: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	if err := b.client.Publish(ctx, broadcastChannel, envelope).Err(); err != nil {
		b.logger.Warn("failed to publish ws message to redis",
			slog.String("error", err.Error()),
			slog.String("user_id", userID),
		)
	}
	return nil
}

// subscribe consumes published messages and delivers each to this instance's
// local sockets via the hub. It runs until Stop is called.
func (b *Broadcaster) subscribe() {
	defer close(b.done)
	ch := b.pubsub.Channel()
	for {
		select {
		case <-b.stop:
			return
		case redisMsg, ok := <-ch:
			if !ok {
				return
			}
			var envelope userMessage
			if err := json.Unmarshal([]byte(redisMsg.Payload), &envelope); err != nil {
				b.logger.Warn("invalid ws envelope from redis", slog.String("error", err.Error()))
				continue
			}
			if err := b.local.SendToUser(envelope.UserID, envelope.Payload); err != nil {
				b.logger.Warn("failed to deliver ws message locally",
					slog.String("error", err.Error()),
					slog.String("user_id", envelope.UserID),
				)
			}
		}
	}
}

// Stop shuts down the subscriber goroutine and closes the pub/sub subscription.
func (b *Broadcaster) Stop() {
	close(b.stop)
	_ = b.pubsub.Close()
	<-b.done
}
