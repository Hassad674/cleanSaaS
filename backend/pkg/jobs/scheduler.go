package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

type Job struct {
	Name     string
	Interval time.Duration
	Fn       func(ctx context.Context) error
}

type Scheduler struct {
	jobs   []Job
	stop   chan struct{}
	wg     sync.WaitGroup
	logger *slog.Logger
}

func NewScheduler(logger *slog.Logger) *Scheduler {
	return &Scheduler{
		stop:   make(chan struct{}),
		logger: logger,
	}
}

func (s *Scheduler) Register(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		s.wg.Add(1)
		go s.runJob(ctx, job)
	}
	s.logger.Info("job scheduler started", slog.Int("jobs", len(s.jobs)))
}

func (s *Scheduler) Stop() {
	close(s.stop)
	s.wg.Wait()
	s.logger.Info("job scheduler stopped")
}

func (s *Scheduler) runJob(ctx context.Context, job Job) {
	defer s.wg.Done()

	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	s.logger.Info("job registered", slog.String("name", job.Name), slog.String("interval", job.Interval.String()))

	for {
		select {
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			if err := s.executeOnce(ctx, job); err != nil {
				s.logger.Error("job failed",
					slog.String("name", job.Name),
					slog.String("error", err.Error()),
					slog.Duration("duration", time.Since(start)),
				)
			} else {
				s.logger.Info("job completed",
					slog.String("name", job.Name),
					slog.Duration("duration", time.Since(start)),
				)
			}
		}
	}
}

// executeOnce runs a single job invocation and recovers from panics, so a single bad
// job tick can never crash the process or permanently kill the scheduler goroutine.
func (s *Scheduler) executeOnce(ctx context.Context, job Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("job panicked (recovered)",
				slog.String("name", job.Name),
				slog.Any("panic", r),
				slog.String("stack", string(debug.Stack())),
			)
			err = fmt.Errorf("job %q panicked: %v", job.Name, r)
		}
	}()
	return job.Fn(ctx)
}
