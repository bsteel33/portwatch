// Package portretry provides retry logic for transient port scan failures.
// It retries a scan operation up to a configured maximum number of attempts
// with an optional backoff delay between each attempt.
package portretry

import (
	"errors"
	"time"
)

// ScanFunc is a function that performs a scan and returns an error on failure.
type ScanFunc func() error

// Retryer retries a scan function on failure.
type Retryer struct {
	cfg Config
	now func() time.Time
	sleep func(time.Duration)
}

// New creates a new Retryer with the given config.
func New(cfg Config) *Retryer {
	return &Retryer{
		cfg:   cfg,
		now:   time.Now,
		sleep: time.Sleep,
	}
}

// Run executes fn up to MaxAttempts times. It returns nil on first success.
// If all attempts fail, it returns the last error wrapped with attempt count.
func (r *Retryer) Run(fn ScanFunc) error {
	var last error
	for i := 0; i < r.cfg.MaxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			last = err
		}
		if i < r.cfg.MaxAttempts-1 && r.cfg.Delay > 0 {
			r.sleep(r.cfg.Delay)
		}
	}
	if last == nil {
		return nil
	}
	return errors.New("all attempts failed: " + last.Error())
}

// Attempts returns the configured maximum number of attempts.
func (r *Retryer) Attempts() int {
	return r.cfg.MaxAttempts
}
