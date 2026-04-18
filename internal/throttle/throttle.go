// Package throttle limits how frequently alerts are emitted for the same port.
package throttle

import (
	"sync"
	"time"
)

// Throttle suppresses repeated alerts for the same key within a cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New returns a Throttle with the given cooldown duration.
func New(cfg Config) *Throttle {
	return &Throttle{
		cooldown: cfg.Cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// It updates the last-seen timestamp when returning true.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	if last, ok := t.last[key]; ok && now.Sub(last) < t.cooldown {
		return false
	}
	t.last[key] = now
	return true
}

// Reset clears the throttle state for a specific key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// ResetAll clears all throttle state.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[string]time.Time)
}
