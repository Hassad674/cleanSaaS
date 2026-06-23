package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"log/slog"
)

// AdvisoryLock implements the scheduler's leader-election Lock using PostgreSQL
// session-level advisory locks (pg_try_advisory_lock / pg_advisory_unlock). It
// needs ZERO extra infrastructure — the database the app already depends on is
// the coordination point — which is why scheduler multi-instance safety does not
// require Redis.
//
// Mechanism: each job name is hashed to a stable 64-bit key. TryAcquire checks
// out a dedicated pooled connection and calls pg_try_advisory_lock(key); the lock
// is held by that SESSION, so the same connection is kept for the lock's lifetime
// and released (advisory-unlock + connection close) by the returned release func.
// Across instances only one session can hold a given key at a time, so only one
// instance runs a given job tick.
type AdvisoryLock struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewAdvisoryLock builds a Postgres advisory-lock leader-election primitive.
func NewAdvisoryLock(db *sql.DB, logger *slog.Logger) *AdvisoryLock {
	if logger == nil {
		logger = slog.Default()
	}
	return &AdvisoryLock{db: db, logger: logger}
}

// TryAcquire attempts to take the advisory lock for name without blocking.
//
//   - (true, release, nil)  — lock taken; caller MUST call release when done.
//   - (false, nil, nil)     — another session/instance holds it; skip this tick.
//   - (false, nil, err)     — the database errored; the scheduler skips the tick.
//
// release is idempotent: it unlocks (best effort) and returns the dedicated
// connection to the pool exactly once.
func (l *AdvisoryLock) TryAcquire(ctx context.Context, name string) (bool, func(), error) {
	key := advisoryLockKey(name)

	conn, err := l.db.Conn(ctx)
	if err != nil {
		return false, nil, fmt.Errorf("acquiring connection for advisory lock: %w", err)
	}

	var acquired bool
	if err := conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", key).Scan(&acquired); err != nil {
		_ = conn.Close()
		return false, nil, fmt.Errorf("pg_try_advisory_lock: %w", err)
	}

	if !acquired {
		// Another session holds the lock — return the connection to the pool.
		_ = conn.Close()
		return false, nil, nil
	}

	var released bool
	release := func() {
		if released {
			return
		}
		released = true
		// Unlock on the SAME session that locked, then close the connection.
		// Use a fresh background context so release still runs if the tick's
		// context was cancelled (e.g. on shutdown).
		if _, err := conn.ExecContext(context.Background(), "SELECT pg_advisory_unlock($1)", key); err != nil {
			l.logger.Warn("failed to release advisory lock",
				slog.String("name", name),
				slog.String("error", err.Error()),
			)
		}
		_ = conn.Close()
	}

	return true, release, nil
}

// advisoryLockKey hashes a job name to a stable signed 64-bit key for
// pg_advisory_lock, which takes a bigint. FNV-1a gives a deterministic value, so
// every instance maps a given job name to the SAME lock key.
func advisoryLockKey(name string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(name))
	return int64(h.Sum64())
}
