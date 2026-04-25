// Package portconn tracks active connection counts per port.
package portconn

import (
	"fmt"
	"sort"
	"sync"
)

// Entry holds connection tracking data for a single port.
type Entry struct {
	Port  int
	Proto string
	Count int
}

// Tracker counts active connections per port/proto pair.
type Tracker struct {
	mu      sync.Mutex
	counts  map[string]Entry
	thresh  int
}

// New creates a Tracker with an optional alert threshold (0 = disabled).
func New(threshold int) *Tracker {
	return &Tracker{
		counts: make(map[string]Entry),
		thresh: threshold,
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Set records the current connection count for a port/proto pair.
func (t *Tracker) Set(port int, proto string, count int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[key(port, proto)] = Entry{Port: port, Proto: proto, Count: count}
}

// Get returns the tracked entry for a port/proto pair and whether it exists.
func (t *Tracker) Get(port int, proto string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.counts[key(port, proto)]
	return e, ok
}

// Exceeded returns all entries whose connection count exceeds the threshold.
// Returns nil if threshold is 0 (disabled).
func (t *Tracker) Exceeded() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.thresh <= 0 {
		return nil
	}
	var out []Entry
	for _, e := range t.counts {
		if e.Count > t.thresh {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Port != out[j].Port {
			return out[i].Port < out[j].Port
		}
		return out[i].Proto < out[j].Proto
	})
	return out
}

// Reset clears all tracked counts.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts = make(map[string]Entry)
}
