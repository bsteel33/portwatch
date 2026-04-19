package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a maximum number of events per time window.
type Limiter struct {
	mu       sync.Mutex
	config   Config
	timestamps []time.Time
	now      func() time.Time
}

// New creates a Limiter with the given config.
func New(cfg Config) *Limiter {
	return &Limiter{
		config: cfg,
		now:    time.Now,
	}
}

// Allow returns true if the event is permitted under the rate limit.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	window := now.Add(-l.config.Window)

	// Evict timestamps outside the window
	filtered := l.timestamps[:0]
	for _, t := range l.timestamps {
		if t.After(window) {
			filtered = append(filtered, t)
		}
	}
	l.timestamps = filtered

	if len(l.timestamps) >= l.config.MaxEvents {
		return false
	}

	l.timestamps = append(l.timestamps, now)
	return true
}

// Reset clears all recorded timestamps.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timestamps = nil
}

// Remaining returns how many events are still allowed in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	window := now.Add(-l.config.Window)
	count := 0
	for _, t := range l.timestamps {
		if t.After(window) {
			count++
		}
	}
	remaining := l.config.MaxEvents - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// RetryAfter returns the duration until at least one additional event is allowed.
// If events are currently allowed, it returns 0.
func (l *Limiter) RetryAfter() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.timestamps) < l.config.MaxEvents {
		return 0
	}

	now := l.now()
	// The oldest timestamp in the window determines when a slot frees up.
	oldest := l.timestamps[0]
	retryAt := oldest.Add(l.config.Window)
	if retryAt.Before(now) {
		return 0
	}
	return retryAt.Sub(now)
}
