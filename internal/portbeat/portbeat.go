// Package portbeat tracks heartbeat (liveness) state for monitored ports,
// recording the last time each port was seen open and detecting stale ports
// that have not been observed within a configurable window.
package portbeat

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Beat holds the last-seen timestamp for a single port.
type Beat struct {
	Port     int
	Proto    string
	LastSeen time.Time
}

// Tracker maintains heartbeat state for observed ports.
type Tracker struct {
	mu      sync.Mutex
	beats   map[string]Beat
	staleIn time.Duration
	now     func() time.Time
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New creates a Tracker that considers a port stale after staleIn duration.
func New(staleIn time.Duration) *Tracker {
	return &Tracker{
		beats:   make(map[string]Beat),
		staleIn: staleIn,
		now:     time.Now,
	}
}

// Pulse records the current time as the last-seen moment for each port in ports.
func (t *Tracker) Pulse(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	for _, p := range ports {
		t.beats[key(p.Port, p.Proto)] = Beat{
			Port:     p.Port,
			Proto:    p.Proto,
			LastSeen: now,
		}
	}
}

// Stale returns all ports whose last heartbeat is older than the stale window.
func (t *Tracker) Stale() []Beat {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := t.now().Add(-t.staleIn)
	var out []Beat
	for _, b := range t.beats {
		if b.LastSeen.Before(cutoff) {
			out = append(out, b)
		}
	}
	return out
}

// Get returns the Beat for a specific port/proto, and whether it exists.
func (t *Tracker) Get(port int, proto string) (Beat, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	b, ok := t.beats[key(port, proto)]
	return b, ok
}

// Reset clears all heartbeat state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.beats = make(map[string]Beat)
}
