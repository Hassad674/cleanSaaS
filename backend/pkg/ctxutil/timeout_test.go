package ctxutil

import (
	"context"
	"testing"
	"time"
)

func TestWithTimeout_AppliesDefaultWhenNoDeadline(t *testing.T) {
	ctx, cancel := WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline to be set on the derived context")
	}
	if remaining := time.Until(dl); remaining <= 0 || remaining > 50*time.Millisecond {
		t.Fatalf("deadline out of expected range: %s", remaining)
	}
}

func TestWithTimeout_KeepsNearerExistingDeadline(t *testing.T) {
	parent, parentCancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer parentCancel()

	parentDeadline, _ := parent.Deadline()

	// Proposing a far larger default must NOT push the deadline out.
	ctx, cancel := WithTimeout(parent, time.Hour)
	defer cancel()

	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected the existing deadline to be preserved")
	}
	if !dl.Equal(parentDeadline) {
		t.Fatalf("nearer parent deadline should win: got %s, want %s", dl, parentDeadline)
	}
}

func TestWithTimeout_TightensFartherExistingDeadline(t *testing.T) {
	parent, parentCancel := context.WithTimeout(context.Background(), time.Hour)
	defer parentCancel()

	ctx, cancel := WithTimeout(parent, 30*time.Millisecond)
	defer cancel()

	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline")
	}
	if remaining := time.Until(dl); remaining > 30*time.Millisecond {
		t.Fatalf("expected the tighter default to apply, remaining=%s", remaining)
	}
}

func TestWithTimeout_NonPositiveMeansNoTimeout(t *testing.T) {
	ctx, cancel := WithTimeout(context.Background(), 0)
	defer cancel()

	if _, ok := ctx.Deadline(); ok {
		t.Fatal("non-positive timeout must not impose a deadline")
	}
}
