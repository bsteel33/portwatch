// Package portcount tracks the number of open ports over time and
// provides threshold-based alerting when the count exceeds a limit.
package portcount

import (
	"fmt"
	"sync"
)

// Config holds portcount configuration.
type Config struct {
	MaxPorts int // alert threshold; 0 means no limit
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{MaxPorts: 0}
}

// Counter tracks open port counts.
type Counter struct {
	mu      sync.Mutex
	cfg     Config
	current int
	peak    int
}

// New creates a new Counter with the given config.
func New(cfg Config) *Counter {
	return &Counter{cfg: cfg}
}

// Update sets the current open port count and returns an alert message
// if the count exceeds the configured maximum. Returns empty string when
// no threshold is set or the count is within limits.
func (c *Counter) Update(n int) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = n
	if n > c.peak {
		c.peak = n
	}
	if c.cfg.MaxPorts > 0 && n > c.cfg.MaxPorts {
		return fmt.Sprintf("open port count %d exceeds threshold %d", n, c.cfg.MaxPorts)
	}
	return ""
}

// Current returns the most recently recorded port count.
func (c *Counter) Current() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}

// Peak returns the highest port count ever recorded.
func (c *Counter) Peak() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.peak
}

// Reset clears current and peak counts.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = 0
	c.peak = 0
}
