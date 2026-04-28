// Package portprobe tracks active probe results per port, recording
// latency and reachability state for each observed endpoint.
package portprobe

import (
	"fmt"
	"sync"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Port      int
	Proto     string
	Reachable bool
	Latency   time.Duration
	CheckedAt time.Time
}

// Tracker stores the most recent probe result for each port/proto pair.
type Tracker struct {
	mu      sync.RWMutex
	results map[string]Result
	clock   func() time.Time
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		results: make(map[string]Result),
		clock:   time.Now,
	}
}

// Record stores a probe result, stamping it with the current time.
func (t *Tracker) Record(port int, proto string, reachable bool, latency time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.results[key(port, proto)] = Result{
		Port:      port,
		Proto:     proto,
		Reachable: reachable,
		Latency:   latency,
		CheckedAt: t.clock(),
	}
}

// Get returns the most recent result for a port/proto pair and whether it exists.
func (t *Tracker) Get(port int, proto string) (Result, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.results[key(port, proto)]
	return r, ok
}

// All returns a snapshot of all stored results.
func (t *Tracker) All() []Result {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Result, 0, len(t.results))
	for _, r := range t.results {
		out = append(out, r)
	}
	return out
}

// Unreachable returns all results where Reachable is false.
func (t *Tracker) Unreachable() []Result {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []Result
	for _, r := range t.results {
		if !r.Reachable {
			out = append(out, r)
		}
	}
	return out
}

// Reset clears all stored results.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.results = make(map[string]Result)
}
