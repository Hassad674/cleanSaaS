package jobs

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// sharedLock is a fake leader-election Lock shared by multiple in-process
// schedulers in tests. It models the real cross-instance guarantee: at most one
// holder of a given name at a time. It lets us prove the scheduler runs a job at
// most once across contending instances WITHOUT needing a real database.
type sharedLock struct {
	mu   sync.Mutex
	held map[string]bool
}

func newSharedLock() *sharedLock {
	return &sharedLock{held: make(map[string]bool)}
}

func (l *sharedLock) TryAcquire(_ context.Context, name string) (bool, func(), error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.held[name] {
		return false, nil, nil
	}
	l.held[name] = true
	release := func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.held[name] = false
	}
	return true, release, nil
}

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestScheduler_NoLock_RunsEveryTick confirms the default (no lock) behavior is
// unchanged: with no leader-election lock, the scheduler runs every tick locally.
func TestScheduler_NoLock_RunsEveryTick(t *testing.T) {
	s := NewScheduler(quietLogger())

	var count int64
	s.Register(Job{
		Name:     "job",
		Interval: 30 * time.Millisecond,
		Fn: func(_ context.Context) error {
			atomic.AddInt64(&count, 1)
			return nil
		},
	})

	s.Start(context.Background())
	time.Sleep(150 * time.Millisecond)
	s.Stop()

	assert.GreaterOrEqual(t, atomic.LoadInt64(&count), int64(2), "no-lock scheduler should run every tick")
}

// TestScheduler_SharedLock_RunsAtMostOnceAcrossInstances proves the core
// multi-instance guarantee: two schedulers sharing one lock and ticking in lock
// step never both run the same job tick. We count total executions across both
// "instances"; with a shared lock a tick is executed by exactly one of them, so
// the combined count never exceeds the number of distinct ticks.
func TestScheduler_SharedLock_RunsAtMostOnceAcrossInstances(t *testing.T) {
	lock := newSharedLock()

	var total int64
	makeScheduler := func() *Scheduler {
		s := NewScheduler(quietLogger())
		s.SetLock(lock)
		s.Register(Job{
			Name:     "shared-job",
			Interval: 30 * time.Millisecond,
			Fn: func(_ context.Context) error {
				// Hold the lock briefly so a contending instance ticking at the
				// same time is forced to skip — mirrors a real job's duration.
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt64(&total, 1)
				return nil
			},
		})
		return s
	}

	a := makeScheduler()
	b := makeScheduler()

	ctx := context.Background()
	a.Start(ctx)
	b.Start(ctx)

	time.Sleep(300 * time.Millisecond)

	a.Stop()
	b.Stop()

	// Over ~300ms at a 30ms interval there are ~10 ticks per instance. Without the
	// lock the two instances would run ~20 times; with the lock each wall-clock
	// tick is run by at most one instance. We assert the combined count stays in
	// the single-instance range (well under 2x), proving dedup across instances.
	runs := atomic.LoadInt64(&total)
	assert.Greater(t, runs, int64(0), "the job should have run at least once")
	assert.LessOrEqual(t, runs, int64(12), "shared lock must prevent both instances running the same tick (got %d)", runs)
}
