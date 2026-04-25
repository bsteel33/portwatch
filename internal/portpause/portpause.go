// Package portpause provides a mechanism to temporarily pause alerting
// for specific ports, suppressing notifications for a configured duration.
package portpause

import (
	"fmt"
	"sync"
	"time"
)

// Entry represents a paused port entry.
type Entry struct {
	Port     int
	Proto    string
	ResumeAt time.Time
}

// Pauser tracks paused ports and their resume times.
type Pauser struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Pauser.
func New() *Pauser {
	return &Pauser{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Pause suppresses alerts for the given port/proto for the specified duration.
func (p *Pauser) Pause(port int, proto string, d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries[key(port, proto)] = Entry{
		Port:     port,
		Proto:    proto,
		ResumeAt: p.now().Add(d),
	}
}

// Resume removes the pause for the given port/proto immediately.
func (p *Pauser) Resume(port int, proto string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.entries, key(port, proto))
}

// IsPaused returns true if the port/proto is currently paused.
func (p *Pauser) IsPaused(port int, proto string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	e, ok := p.entries[key(port, proto)]
	if !ok {
		return false
	}
	if p.now().After(e.ResumeAt) {
		delete(p.entries, key(port, proto))
		return false
	}
	return true
}

// Active returns all currently paused entries.
func (p *Pauser) Active() []Entry {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	var out []Entry
	for k, e := range p.entries {
		if now.After(e.ResumeAt) {
			delete(p.entries, k)
			continue
		}
		out = append(out, e)
	}
	return out
}

// Reset clears all paused entries.
func (p *Pauser) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries = make(map[string]Entry)
}
