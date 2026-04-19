// Package portburst detects sudden bursts in newly opened ports within a time window.
package portburst

import (
	"sync"
	"time"
)

// Event represents a burst detection result.
type Event struct {
	Count     int
	Threshold int
	Window    time.Duration
	At        time.Time
}

// Detector tracks port open events and fires when a burst threshold is exceeded.
type Detector struct {
	mu     sync.Mutex
	cfg    Config
	events []time.Time
	clock  func() time.Time
}

// New returns a new Detector with the given config.
func New(cfg Config) *Detector {
	return &Detector{cfg: cfg, clock: time.Now}
}

// Record registers a batch of newly opened port count and returns a burst Event
// if the threshold is exceeded within the window, or nil otherwise.
func (d *Detector) Record(count int) *Event {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	for i := 0; i < count; i++ {
		d.events = append(d.events, now)
	}
	d.prune(now)

	if len(d.events) >= d.cfg.Threshold {
		return &Event{
			Count:     len(d.events),
			Threshold: d.cfg.Threshold,
			Window:    d.cfg.Window,
			At:        now,
		}
	}
	return nil
}

// Reset clears all recorded events.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = nil
}

func (d *Detector) prune(now time.Time) {
	cutoff := now.Add(-d.cfg.Window)
	i := 0
	for i < len(d.events) && d.events[i].Before(cutoff) {
		i++
	}
	d.events = d.events[i:]
}
