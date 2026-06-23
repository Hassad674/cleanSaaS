package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
)

// Limiter is the rate-limiting abstraction the middleware depends on. It lets
// the app swap a single-instance in-memory limiter for a distributed
// (Redis-backed) one without the HTTP layer knowing which is in use. Allow
// reports whether the given key (client IP) may proceed; Stop releases any
// background resources held by the implementation.
type Limiter interface {
	Allow(key string) bool
	Stop()
}

// NewLimiter builds a Limiter allowing requestsPerMinute per client IP, in a
// distinct keyspace (so independent limiters never share counters). The
// composition root injects this so the router does not need to know whether the
// backing store is in-memory or Redis. Each call to a NewLimiter implementation
// returns a fresh limiter the router owns for the lifetime of its routes.
type NewLimiter func(requestsPerMinute float64, keyspace string) Limiter

// InMemoryLimiterFactory returns a NewLimiter that always builds the in-process
// token-bucket limiter, ignoring the keyspace (each limiter holds its own map).
// This is the default and the fail-open fallback when no Redis is configured.
func InMemoryLimiterFactory() NewLimiter {
	return func(requestsPerMinute float64, _ string) Limiter {
		return NewRateLimiter(requestsPerMinute)
	}
}

// bucket holds the token bucket state for one IP.
type bucket struct {
	tokens    float64
	lastCheck time.Time
}

// RateLimiter is the in-process, single-instance token-bucket limiter (per IP).
// It is the default and the fail-open fallback when no Redis is configured.
// Across multiple instances its limits multiply per instance — use the
// Redis-backed Limiter for horizontally-scaled deployments.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64 // tokens per second
	capacity float64 // max tokens
	stop     chan struct{}
}

// NewRateLimiter creates a new rate limiter.
// rate is requests per minute, capacity equals rate.
func NewRateLimiter(requestsPerMinute float64) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     requestsPerMinute / 60.0,
		capacity: requestsPerMinute,
		stop:     make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

// Stop signals the cleanup goroutine to exit.
func (rl *RateLimiter) Stop() {
	close(rl.stop)
}

// Allow checks if the IP is allowed to make a request.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[ip]
	if !exists {
		rl.buckets[ip] = &bucket{
			tokens:    rl.capacity - 1,
			lastCheck: now,
		}
		return true
	}

	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.capacity {
		b.tokens = rl.capacity
	}
	b.lastCheck = now

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

// cleanup removes stale entries every 5 minutes. Stops when Stop() is called.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stop:
			return
		case <-ticker.C:
			rl.mu.Lock()
			cutoff := time.Now().Add(-10 * time.Minute)
			for ip, b := range rl.buckets {
				if b.lastCheck.Before(cutoff) {
					delete(rl.buckets, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// rateLimitExemptPaths are operational endpoints that must never be throttled:
// liveness/readiness probes and the Prometheus scrape target are polled
// frequently by infrastructure and rejecting them would cause false outages.
var rateLimitExemptPaths = map[string]bool{
	"/health":  true,
	"/livez":   true,
	"/readyz":  true,
	"/metrics": true,
}

// RateLimit returns middleware that rate-limits by client IP, exempting the
// operational endpoints in rateLimitExemptPaths. It depends on the Limiter
// abstraction, so the same middleware works with either the in-memory limiter
// (single instance) or a Redis-backed one (shared across instances).
func RateLimit(limiter Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rateLimitExemptPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}
			if !limiter.Allow(r.RemoteAddr) {
				w.Header().Set("Retry-After", "60")
				response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
