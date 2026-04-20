// Package portcache provides a short-lived in-memory cache for scanned port
// results, reducing redundant scans within a configurable TTL window.
package portcache

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// entry holds a cached scan result with its expiry time.
type entry struct {
	ports  []scanner.Port
	expiry time.Time
}

// Cache stores the most recent scan result and serves it until the TTL expires.
type Cache struct {
	mu    sync.Mutex
	data  *entry
	ttl   time.Duration
	clock func() time.Time
}

// New creates a Cache with the given TTL.
func New(cfg Config) *Cache {
	return &Cache{
		ttl:   cfg.TTL,
		clock: time.Now,
	}
}

// Set stores ports in the cache, replacing any previous value.
func (c *Cache) Set(ports []scanner.Port) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = &entry{
		ports:  ports,
		expiry: c.clock().Add(c.ttl),
	}
}

// Get returns the cached ports and true if the entry exists and has not expired.
// Returns nil and false otherwise.
func (c *Cache) Get() ([]scanner.Port, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil {
		return nil, false
	}
	if c.clock().After(c.data.expiry) {
		c.data = nil
		return nil, false
	}
	return c.data.ports, true
}

// Invalidate clears the cached entry immediately.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = nil
}
