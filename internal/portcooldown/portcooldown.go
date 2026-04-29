// Package portcooldown tracks per-port cooldown periods after state changes,
// preventing repeated alerts for the same port within a configurable window.
package portcooldown

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds cooldown state for a single port.
type Entry struct {
	Port     int
	Proto    string
	CooledAt time.Time
}

// Cooldown manages per-port cooldown tracking.
type Cooldown struct {
	mu      sync.Mutex
	entries map[string]Entry
	window  time.Duration
	now     func() time.Time
}

// New returns a Cooldown with the given window duration.
func New(window time.Duration) *Cooldown {
	return &Cooldown{
		entries: make(map[string]Entry),
		window:  window,
		now:     time.Now,
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Mark records the current time as the start of a cooldown for the given port.
func (c *Cooldown) Mark(port int, proto string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key(port, proto)] = Entry{
		Port:     port,
		Proto:    proto,
		CooledAt: c.now(),
	}
}

// IsCooling returns true if the port is within its cooldown window.
func (c *Cooldown) IsCooling(port int, proto string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key(port, proto)]
	if !ok {
		return false
	}
	return c.now().Before(e.CooledAt.Add(c.window))
}

// Reset clears the cooldown entry for a port.
func (c *Cooldown) Reset(port int, proto string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key(port, proto))
}

// Active returns all entries still within their cooldown window.
func (c *Cooldown) Active() []Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	var out []Entry
	for _, e := range c.entries {
		if now.Before(e.CooledAt.Add(c.window)) {
			out = append(out, e)
		}
	}
	return out
}
