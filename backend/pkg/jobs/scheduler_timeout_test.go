package jobs

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

func quietSchedulerWithTimeout(d time.Duration) *Scheduler {
	return NewSchedulerWithTimeout(slog.New(slog.NewTextHandler(io.Discard, nil)), d)
}

func TestScheduler_executeOnce_GivesJobADeadline(t *testing.T) {
	s := quietSchedulerWithTimeout(5 * time.Second)

	var (
		hadDeadline bool
		remaining   time.Duration
	)
	job := Job{Name: "deadline-check", Fn: func(ctx context.Context) error {
		dl, ok := ctx.Deadline()
		hadDeadline = ok
		if ok {
			remaining = time.Until(dl)
		}
		return nil
	}}

	if err := s.executeOnce(context.Background(), job); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hadDeadline {
		t.Fatal("expected the job context to carry a deadline")
	}
	if remaining <= 0 || remaining > 5*time.Second {
		t.Fatalf("deadline out of expected range: %s", remaining)
	}
}

func TestScheduler_executeOnce_PerJobTimeoutOverridesDefault(t *testing.T) {
	s := quietSchedulerWithTimeout(time.Hour)

	var remaining time.Duration
	job := Job{
		Name:    "short-job",
		Timeout: 40 * time.Millisecond,
		Fn: func(ctx context.Context) error {
			dl, ok := ctx.Deadline()
			if !ok {
				t.Fatal("expected a deadline")
			}
			remaining = time.Until(dl)
			return nil
		},
	}

	if err := s.executeOnce(context.Background(), job); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if remaining > 40*time.Millisecond {
		t.Fatalf("per-job timeout should win over the scheduler default, remaining=%s", remaining)
	}
}

func TestScheduler_executeOnce_NegativeTimeoutDisablesDeadline(t *testing.T) {
	s := quietSchedulerWithTimeout(5 * time.Second)

	var hadDeadline bool
	job := Job{
		Name:    "unbounded-job",
		Timeout: -1, // opt out of any deadline
		Fn: func(ctx context.Context) error {
			_, hadDeadline = ctx.Deadline()
			return nil
		},
	}

	if err := s.executeOnce(context.Background(), job); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hadDeadline {
		t.Fatal("a negative per-job timeout must disable the deadline")
	}
}

func TestScheduler_executeOnce_StuckJobIsCancelled(t *testing.T) {
	s := quietSchedulerWithTimeout(30 * time.Millisecond)

	job := Job{Name: "stuck", Fn: func(ctx context.Context) error {
		<-ctx.Done() // simulate a hung operation that respects cancellation
		return ctx.Err()
	}}

	start := time.Now()
	err := s.executeOnce(context.Background(), job)
	if err == nil {
		t.Fatal("expected the stuck job to be cancelled by the deadline")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("stuck job ran far longer than its timeout: %s", elapsed)
	}
}
