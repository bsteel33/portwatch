// Package portpool tracks a named pool of ports and enforces capacity limits.
package portpool

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Pool holds a named set of ports with an optional capacity cap.
type Pool struct {
	mu       sync.RWMutex
	name     string
	capacity int
	ports    map[string]scanner.Port
}

// New creates a Pool with the given name and capacity.
// A capacity of 0 means unlimited.
func New(name string, capacity int) *Pool {
	return &Pool{
		name:     name,
		capacity: capacity,
		ports:    make(map[string]scanner.Port),
	}
}

// Add inserts a port into the pool. Returns an error if the pool is at capacity.
func (p *Pool) Add(port scanner.Port) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	k := key(port)
	if _, exists := p.ports[k]; exists {
		return nil
	}
	if p.capacity > 0 && len(p.ports) >= p.capacity {
		return fmt.Errorf("portpool %q: capacity %d reached", p.name, p.capacity)
	}
	p.ports[k] = port
	return nil
}

// Remove deletes a port from the pool.
func (p *Pool) Remove(port scanner.Port) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.ports, key(port))
}

// Contains reports whether the port is in the pool.
func (p *Pool) Contains(port scanner.Port) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.ports[key(port)]
	return ok
}

// All returns a snapshot of all ports currently in the pool.
func (p *Pool) All() []scanner.Port {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]scanner.Port, 0, len(p.ports))
	for _, v := range p.ports {
		out = append(out, v)
	}
	return out
}

// Len returns the current number of ports in the pool.
func (p *Pool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.ports)
}

// Name returns the pool name.
func (p *Pool) Name() string { return p.name }

// Capacity returns the pool capacity limit (0 = unlimited).
func (p *Pool) Capacity() int { return p.capacity }

func key(port scanner.Port) string {
	return fmt.Sprintf("%d/%s", port.Port, port.Proto)
}
