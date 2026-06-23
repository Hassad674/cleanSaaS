// Package redis provides Redis-backed adapters for the application's shared,
// multi-instance-safe components (distributed rate limiting and cross-instance
// WebSocket fan-out). The go-redis library is imported ONLY here and in this
// package's sibling files — the app/domain core depends on interfaces, never on
// this concrete client, in keeping with the hexagonal dependency rule.
package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// connectTimeout bounds the startup PING so an unreachable Redis can never hang
// boot — callers fall back to in-memory behavior when this fails.
const connectTimeout = 3 * time.Second

// Connect parses a redis:// URL, opens a client, and verifies reachability with a
// bounded PING. It returns an error (rather than crashing) when Redis is
// unreachable so the caller can fail OPEN to the single-instance fallback.
func Connect(url string) (*goredis.Client, error) {
	opts, err := goredis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parsing REDIS_URL: %w", err)
	}

	client := goredis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("pinging redis: %w", err)
	}

	return client, nil
}
