package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/pkg/ctxutil"
)

// defaultJobTimeout bounds a single job invocation when neither the Job nor the
// Scheduler specifies one, so a stuck cleanup can never run forever.
const defaultJobTimeout = 30 * time.Second

// Lock is a best-effort, named, mutually-exclusive lock acquired for the duration
// of a single job tick. It is the leader-election primitive that makes the
// scheduler multi-instance safe: across N instances, only the one that wins the
// lock for a given tick runs that job, so a job executes AT MOST ONCE per tick
// cluster-wide instead of N times.
//
// TryAcquire attempts to take the named lock without blocking. It returns
// (true, release, nil) when the lock was taken — the caller MUST call release
// (idempotent) once the tick completes — or (false, nil, nil) when another
// instance holds it (this tick should be skipped). A non-nil error means the
// lock backend itself failed; the scheduler treats that as "do not run this
// tick" so a flaky backend can never cause duplicate execution.
type Lock interface {
	TryAcquire(ctx context.Context, name string) (acquired bool, release func(), err error)
}

type Job struct {
	Name     string
	Interval time.Duration
	Fn       func(ctx context.Context) error

	// Timeout bounds a single invocation of this job. When zero, the scheduler's
	// default (jobTimeout) applies. A negative value disables the timeout for a
	// genuinely long-running job (use with care).
	Timeout time.Duration
}

type Scheduler struct {
	jobs       []Job
	stop       chan struct{}
	wg         sync.WaitGroup
	logger     *slog.Logger
	jobTimeout time.Duration

	// lock is the optional leader-election primitive. When nil (the default), the
	// scheduler runs every tick locally — correct for a single instance and the
	// out-of-the-box experience. When set (e.g. a Postgres advisory lock), each
	// tick first tries to acquire the per-job lock and SKIPS the tick if another
	// instance holds it, so the job runs at most once across the cluster.
	lock Lock
}

func NewScheduler(logger *slog.Logger) *Scheduler {
	return NewSchedulerWithTimeout(logger, defaultJobTimeout)
}

// NewSchedulerWithTimeout builds a Scheduler whose jobs are each bounded by
// jobTimeout unless a Job overrides it via its own Timeout field.
func NewSchedulerWithTimeout(logger *slog.Logger, jobTimeout time.Duration) *Scheduler {
	return &Scheduler{
		stop:       make(chan struct{}),
		logger:     logger,
		jobTimeout: jobTimeout,
	}
}

// jobTimeoutFor resolves the effective timeout for a job: an explicit per-job
// Timeout wins, otherwise the scheduler default. A negative value means "no
// timeout" and is passed through so the job context carries no deadline.
func (s *Scheduler) jobTimeoutFor(job Job) time.Duration {
	if job.Timeout != 0 {
		return job.Timeout
	}
	return s.jobTimeout
}

// SetLock installs an optional leader-election lock so the scheduler is safe to
// run on multiple instances: each job tick runs at most once cluster-wide. Pass
// nil (or never call this) to keep the single-instance behavior. Call before
// Start.
func (s *Scheduler) SetLock(lock Lock) {
	s.lock = lock
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
			s.tick(ctx, job)
		}
	}
}

// tick runs one scheduled invocation of job. When a leader-election lock is
// installed, it first tries to acquire the per-job lock and SKIPS this tick if
// another instance holds it (or the lock backend errored) — that is what makes
// the scheduler safe across multiple instances. Without a lock it always runs.
func (s *Scheduler) tick(ctx context.Context, job Job) {
	if s.lock != nil {
		acquired, release, err := s.lock.TryAcquire(ctx, job.Name)
		if err != nil {
			s.logger.Warn("job tick skipped: lock acquisition failed",
				slog.String("name", job.Name),
				slog.String("error", err.Error()),
			)
			return
		}
		if !acquired {
			s.logger.Debug("job tick skipped: another instance holds the lock",
				slog.String("name", job.Name),
			)
			return
		}
		defer release()
	}

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

// executeOnce runs a single job invocation with a bounded context and recovers
// from panics, so a single bad job tick can neither run forever nor crash the
// process / permanently kill the scheduler goroutine. The per-job deadline is
// derived via ctxutil.WithTimeout (a ceiling — a nearer parent deadline wins).
func (s *Scheduler) executeOnce(ctx context.Context, job Job) (err error) {
	ctx, cancel := ctxutil.WithTimeout(ctx, s.jobTimeoutFor(job))
	defer cancel()

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
