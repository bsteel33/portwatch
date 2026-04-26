// Package portmute provides per-port muting so that known-noisy ports
// can be silenced for a configurable duration without affecting other alerts.
package portmute

import (
	"sync"
	"time"
)

// Entry holds mute state for a single port key.
type Entry struct {
	Until  time.Time
	Reason string
}

// Muter tracks muted ports.
type Muter struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Muter.
func New() *Muter {
	return &Muter{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 8)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

// Mute silences port/proto for the given duration with an optional reason.
func (m *Muter) Mute(port int, proto string, d time.Duration, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key(port, proto)] = Entry{
		Until:  m.now().Add(d),
		Reason: reason,
	}
}

// Unmute removes a mute entry immediately.
func (m *Muter) Unmute(port int, proto string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key(port, proto))
}

// IsMuted reports whether port/proto is currently muted.
func (m *Muter) IsMuted(port int, proto string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[key(port, proto)]
	if !ok {
		return false
	}
	if m.now().After(e.Until) {
		delete(m.entries, key(port, proto))
		return false
	}
	return true
}

// Get returns the mute entry for port/proto, and whether it exists and is active.
func (m *Muter) Get(port int, proto string) (Entry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[key(port, proto)]
	if !ok {
		return Entry{}, false
	}
	if m.now().After(e.Until) {
		delete(m.entries, key(port, proto))
		return Entry{}, false
	}
	return e, true
}

// Reset clears all mute entries.
func (m *Muter) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = make(map[string]Entry)
}
