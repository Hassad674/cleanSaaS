package jobs

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

func quietScheduler() *Scheduler {
	return NewScheduler(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestScheduler_executeOnce_RecoversPanic(t *testing.T) {
	s := quietScheduler()
	job := Job{Name: "boom", Fn: func(context.Context) error { panic("kaboom") }}

	err := s.executeOnce(context.Background(), job)
	if err == nil {
		t.Fatal("expected a recovered panic to be returned as an error")
	}
}

func TestScheduler_executeOnce_PassesThroughResult(t *testing.T) {
	s := quietScheduler()

	if err := s.executeOnce(context.Background(), Job{Name: "ok", Fn: func(context.Context) error { return nil }}); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	sentinel := errors.New("expected failure")
	if err := s.executeOnce(context.Background(), Job{Name: "fail", Fn: func(context.Context) error { return sentinel }}); !errors.Is(err, sentinel) {
		t.Errorf("expected the job's error to pass through, got %v", err)
	}
}
