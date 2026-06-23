package redis

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// rateLimitWindow is the fixed window over which requests are counted. The
// middleware expresses limits as "requests per minute", so the window is 1m.
const rateLimitWindow = time.Minute

// rateLimitTimeout bounds a single Allow() Redis round-trip so a slow/unreachable
// Redis can never block an HTTP request — on timeout we fail OPEN (allow).
const rateLimitTimeout = 200 * time.Millisecond

// slidingWindowScript atomically increments the request counter for a key and
// (only on the first request of a window) sets the window TTL, then returns the
// current count. Doing both in one Lua script makes the check-and-set atomic
// across all instances sharing this Redis, so the limit is GLOBAL rather than
// per-instance. KEYS[1] = counter key, ARGV[1] = window TTL in milliseconds.
var slidingWindowScript = goredis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
return current
`)

// RateLimiter is a Redis-backed fixed-window rate limiter. Because the counter
// lives in Redis, the limit is enforced across every instance sharing the same
// Redis — solving the "limits multiply per instance" problem of the in-memory
// limiter. It satisfies the middleware.Limiter interface.
type RateLimiter struct {
	client   *goredis.Client
	limit    int64
	keyspace string
	logger   *slog.Logger
}

// NewRateLimiter builds a Redis-backed limiter allowing requestsPerMinute per key
// (client IP). keyspace namespaces the Redis keys so multiple limiters (api, auth,
// demo, ...) sharing one Redis never collide.
func NewRateLimiter(client *goredis.Client, requestsPerMinute float64, keyspace string, logger *slog.Logger) *RateLimiter {
	if logger == nil {
		logger = slog.Default()
	}
	return &RateLimiter{
		client:   client,
		limit:    int64(requestsPerMinute),
		keyspace: keyspace,
		logger:   logger,
	}
}

// Allow reports whether the key may proceed. It increments the shared counter and
// allows while the count is within the limit for the current window. On any Redis
// error (timeout/unreachable) it fails OPEN (allows) and logs, so a cache outage
// degrades to "no limiting" rather than rejecting all traffic.
func (rl *RateLimiter) Allow(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), rateLimitTimeout)
	defer cancel()

	redisKey := "ratelimit:" + rl.keyspace + ":" + key
	ttlMillis := strconv.FormatInt(rateLimitWindow.Milliseconds(), 10)

	count, err := slidingWindowScript.Run(ctx, rl.client, []string{redisKey}, ttlMillis).Int64()
	if err != nil {
		rl.logger.Warn("redis rate limiter error, failing open",
			slog.String("error", err.Error()),
			slog.String("keyspace", rl.keyspace),
		)
		return true
	}

	return count <= rl.limit
}

// Stop is a no-op: the Redis client is shared and owned by the composition root,
// which closes it on shutdown. It exists to satisfy the middleware.Limiter
// interface alongside the in-memory limiter.
func (rl *RateLimiter) Stop() {}
