// Package portdrain tracks ports that are in a "draining" state —
// open but scheduled for shutdown — and reports when they finally close.
package portdrain

import (
	"fmt"
	"sync"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

// Entry records when a port was marked for draining.
type Entry struct {
	Port     scanner.Port
	MarkedAt time.Time
	Deadline time.Time
}

// Drainer tracks ports that are expected to close within a deadline.
type Drainer struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

func key(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Port, p.Proto)
}

// New returns a new Drainer.
func New() *Drainer {
	return &Drainer{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Mark registers a port as draining with the given TTL.
func (d *Drainer) Mark(p scanner.Port, ttl time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	d.entries[key(p)] = Entry{
		Port:     p,
		MarkedAt: now,
		Deadline: now.Add(ttl),
	}
}

// IsDraining reports whether the port is currently marked for draining.
func (d *Drainer) IsDraining(p scanner.Port) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	e, ok := d.entries[key(p)]
	if !ok {
		return false
	}
	return d.now().Before(e.Deadline)
}

// Overdue returns all ports whose drain deadline has passed but are still open.
func (d *Drainer) Overdue(open []scanner.Port) []Entry {
	d.mu.Lock()
	defer d.mu.Unlock()
	openSet := make(map[string]struct{}, len(open))
	for _, p := range open {
		openSet[key(p)] = struct{}{}
	}
	now := d.now()
	var out []Entry
	for k, e := range d.entries {
		if _, stillOpen := openSet[k]; stillOpen && now.After(e.Deadline) {
			out = append(out, e)
		}
	}
	return out
}

// Evict removes a port from the drain list (e.g. once it has closed).
func (d *Drainer) Evict(p scanner.Port) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, key(p))
}

// All returns a snapshot of all current drain entries.
func (d *Drainer) All() []Entry {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]Entry, 0, len(d.entries))
	for _, e := range d.entries {
		out = append(out, e)
	}
	return out
}
