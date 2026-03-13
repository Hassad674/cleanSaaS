package jobs

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_StartStop(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := NewScheduler(logger)

	var count int64
	s.Register(Job{
		Name:     "test-job",
		Interval: 50 * time.Millisecond,
		Fn: func(_ context.Context) error {
			atomic.AddInt64(&count, 1)
			return nil
		},
	})

	ctx := context.Background()
	s.Start(ctx)

	time.Sleep(200 * time.Millisecond)
	s.Stop()

	executions := atomic.LoadInt64(&count)
	assert.GreaterOrEqual(t, executions, int64(2), "job should have executed at least twice")
}

func TestScheduler_MultipleJobs(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := NewScheduler(logger)

	var countA, countB int64
	s.Register(Job{
		Name:     "job-a",
		Interval: 50 * time.Millisecond,
		Fn: func(_ context.Context) error {
			atomic.AddInt64(&countA, 1)
			return nil
		},
	})
	s.Register(Job{
		Name:     "job-b",
		Interval: 50 * time.Millisecond,
		Fn: func(_ context.Context) error {
			atomic.AddInt64(&countB, 1)
			return nil
		},
	})

	ctx := context.Background()
	s.Start(ctx)

	time.Sleep(200 * time.Millisecond)
	s.Stop()

	assert.GreaterOrEqual(t, atomic.LoadInt64(&countA), int64(2))
	assert.GreaterOrEqual(t, atomic.LoadInt64(&countB), int64(2))
}

func TestScheduler_StopImmediately(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := NewScheduler(logger)

	s.Register(Job{
		Name:     "long-interval-job",
		Interval: 1 * time.Hour,
		Fn: func(_ context.Context) error {
			return nil
		},
	})

	ctx := context.Background()
	s.Start(ctx)

	// Stop should return quickly
	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()

	select {
	case <-done:
		// OK
	case <-time.After(2 * time.Second):
		t.Fatal("Stop should have returned quickly")
	}
}
