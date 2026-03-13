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
