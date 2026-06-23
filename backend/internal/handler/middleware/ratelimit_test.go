package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow_UnderLimit(t *testing.T) {
	rl := NewRateLimiter(10) // 10 req/min
	for i := 0; i < 10; i++ {
		assert.True(t, rl.Allow("192.168.1.1"), "request %d should be allowed", i)
	}
}

func TestRateLimiter_Allow_OverLimit(t *testing.T) {
	rl := NewRateLimiter(5) // 5 req/min
	for i := 0; i < 5; i++ {
		rl.Allow("10.0.0.1")
	}
	assert.False(t, rl.Allow("10.0.0.1"), "should be rate limited")
}

func TestRateLimiter_Allow_DifferentIPs(t *testing.T) {
	rl := NewRateLimiter(2)
	rl.Allow("10.0.0.1")
	rl.Allow("10.0.0.1")
	assert.False(t, rl.Allow("10.0.0.1"), "IP1 should be limited")
	assert.True(t, rl.Allow("10.0.0.2"), "IP2 should not be limited")
}

// RateLimiter must satisfy the Limiter abstraction the middleware depends on, so
// it stays interchangeable with the Redis-backed limiter.
var _ Limiter = (*RateLimiter)(nil)

func TestInMemoryLimiterFactory_BuildsWorkingLimiter(t *testing.T) {
	newLimiter := InMemoryLimiterFactory()
	lim := newLimiter(2, "ignored-keyspace")
	t.Cleanup(lim.Stop)

	assert.True(t, lim.Allow("10.0.0.9"))
	assert.True(t, lim.Allow("10.0.0.9"))
	assert.False(t, lim.Allow("10.0.0.9"), "third request over the 2/min budget should be limited")
}

func TestInMemoryLimiterFactory_IndependentKeyspaces(t *testing.T) {
	newLimiter := InMemoryLimiterFactory()
	a := newLimiter(1, "ks-a")
	b := newLimiter(1, "ks-b")
	t.Cleanup(a.Stop)
	t.Cleanup(b.Stop)

	// Each factory call returns an independent limiter with its own budget.
	assert.True(t, a.Allow("ip"))
	assert.False(t, a.Allow("ip"))
	assert.True(t, b.Allow("ip"), "a separate limiter must have its own budget")
}
