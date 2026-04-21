// Package portdebounce suppresses transient port changes that appear and
// disappear within a short observation window, reducing alert noise from
// short-lived ephemeral ports.
package portdebounce

import (
	"fmt"
	"sync"
	"time"
)

// Event represents a port open/close observation.
type Event struct {
	Port  int
	Proto string
	Seen  time.Time
}

// Debouncer holds pending port events and confirms them only after a
// stabilisation window has elapsed without a contradicting event.
type Debouncer struct {
	mu      sync.Mutex
	cfg     Config
	now     func() time.Time
	pending map[string]Event // key -> first-seen event
}

// New returns a Debouncer using cfg.
func New(cfg Config) *Debouncer {
	return &Debouncer{
		cfg:     cfg,
		now:     time.Now,
		pending: make(map[string]Event),
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Observe records that a port was seen open. Returns true when the port has
// been stable for at least the configured window and should be acted upon.
func (d *Debouncer) Observe(port int, proto string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	k := key(port, proto)
	now := d.now()

	ev, exists := d.pending[k]
	if !exists {
		d.pending[k] = Event{Port: port, Proto: proto, Seen: now}
		return false
	}

	if now.Sub(ev.Seen) >= d.cfg.Window {
		delete(d.pending, k)
		return true
	}
	return false
}

// Dismiss removes a pending entry — call this when the port is seen closed
// before the window elapses, indicating a transient event.
func (d *Debouncer) Dismiss(port int, proto string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.pending, key(port, proto))
}

// Pending returns the number of events currently waiting in the window.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}

// Reset clears all pending events.
func (d *Debouncer) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pending = make(map[string]Event)
}
