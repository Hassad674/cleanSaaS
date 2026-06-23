//go:build integration

package postgres

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hassad/boilerplateSaaS/backend/pkg/jobs"
)

func itestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestAdvisoryLock_MutualExclusion proves the core leader-election invariant
// against the LIVE database: while one session holds the advisory lock for a
// name, a second contending session cannot acquire it; after release it can.
func TestAdvisoryLock_MutualExclusion(t *testing.T) {
	db := openTestDB(t)
	lock := NewAdvisoryLock(db, itestLogger())

	ctx := context.Background()
	name := "itest-advisory-" + uniqueTag()

	acquired, release, err := lock.TryAcquire(ctx, name)
	require.NoError(t, err)
	require.True(t, acquired, "first acquisition should succeed")

	// A second contender must be refused while the first holds the lock.
	acquired2, release2, err := lock.TryAcquire(ctx, name)
	require.NoError(t, err)
	assert.False(t, acquired2, "second acquisition must be refused while held")
	assert.Nil(t, release2)

	// After release, the lock becomes available again.
	release()
	acquired3, release3, err := lock.TryAcquire(ctx, name)
	require.NoError(t, err)
	assert.True(t, acquired3, "lock should be re-acquirable after release")
	release3()
}

// TestAdvisoryLock_ReleaseIsIdempotent guards the contract that the returned
// release is safe to call more than once (defer + explicit, shutdown paths, ...).
func TestAdvisoryLock_ReleaseIsIdempotent(t *testing.T) {
	db := openTestDB(t)
	lock := NewAdvisoryLock(db, itestLogger())
	name := "itest-advisory-idem-" + uniqueTag()

	acquired, release, err := lock.TryAcquire(context.Background(), name)
	require.NoError(t, err)
	require.True(t, acquired)

	release()
	release() // must not panic or error
}

// TestScheduler_WithAdvisoryLock_RunsOnceAcrossInstances is the end-to-end proof
// requested: two schedulers (simulating two instances) each backed by a Postgres
// advisory lock and registering the SAME job must, between them, run that job at
// most once per tick — never both. We count combined executions and assert they
// stay in the single-instance range.
func TestScheduler_WithAdvisoryLock_RunsOnceAcrossInstances(t *testing.T) {
	db := openTestDB(t)

	var total int64
	jobName := "itest-scheduled-" + uniqueTag()

	makeScheduler := func() *jobs.Scheduler {
		s := jobs.NewScheduler(itestLogger())
		s.SetLock(NewAdvisoryLock(db, itestLogger()))
		s.Register(jobs.Job{
			Name:     jobName,
			Interval: 40 * time.Millisecond,
			Fn: func(_ context.Context) error {
				// Hold long enough that a contending instance ticking concurrently
				// is forced to skip (it cannot acquire the held lock).
				time.Sleep(15 * time.Millisecond)
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

	time.Sleep(400 * time.Millisecond)

	a.Stop()
	b.Stop()

	runs := atomic.LoadInt64(&total)
	assert.Greater(t, runs, int64(0), "the job should have run at least once")
	// ~400ms / 40ms = ~10 ticks per instance. Without the lock the pair would run
	// ~20 times; with the advisory lock each tick is run by exactly one instance.
	assert.LessOrEqual(t, runs, int64(13),
		"advisory lock must prevent both instances running the same tick (got %d)", runs)
}
