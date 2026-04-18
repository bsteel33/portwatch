package schedule

import (
	"context"
	"time"
)

// Schedule drives periodic execution of a job function.
type Schedule struct {
	cfg Config
}

// New returns a Schedule with the given config.
func New(cfg Config) *Schedule {
	return &Schedule{cfg: cfg}
}

// Run calls job on every tick until ctx is cancelled.
// The first call happens immediately unless cfg.DelayFirst is true.
func (s *Schedule) Run(ctx context.Context, job func()) error {
	if !s.cfg.DelayFirst {
		job()
	}
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			job()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// RunN calls job up to n times (or until ctx is cancelled).
func (s *Schedule) RunN(ctx context.Context, n int, job func()) error {
	count := 0
	return s.Run(ctx, func() {
		if count >= n {
			return
		}
		job()
		count++
	})
}
