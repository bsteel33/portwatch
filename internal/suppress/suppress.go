// Package suppress provides a mechanism to suppress repeated alerts
// for ports that are already known to be in a changed state.
package suppress

import (
	"sync"
	"time"
)

// Entry tracks when a port key was first suppressed.
type Entry struct {
	SuppressedAt time.Time
	Expires      time.Time
}

// Suppressor holds suppressed port keys.
type Suppressor struct {
	mu      sync.Mutex
	entries map[string]Entry
	ttl     time.Duration
	now     func() time.Time
}

// New returns a new Suppressor with the given TTL.
func New(cfg Config) *Suppressor {
	return &Suppressor{
		entries: make(map[string]Entry),
		ttl:     cfg.TTL,
		now:     time.Now,
	}
}

// IsSuppressed returns true if the key is currently suppressed.
func (s *Suppressor) IsSuppressed(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	if s.now().After(e.Expires) {
		delete(s.entries, key)
		return false
	}
	return true
}

// Suppress marks the key as suppressed for the configured TTL.
func (s *Suppressor) Suppress(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	s.entries[key] = Entry{
		SuppressedAt: now,
		Expires:      now.Add(s.ttl),
	}
}

// Reset removes the key from suppression.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Len returns the number of currently suppressed keys.
func (s *Suppressor) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}
