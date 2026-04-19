// Package portstate tracks the up/down state of individual ports across scans.
package portstate

import (
	"sync"
	"time"
)

// State represents the observed state of a single port.
type State struct {
	Port      int
	Proto     string
	Up        bool
	LastSeen  time.Time
	FirstSeen time.Time
	Flaps     int
}

// Tracker maintains state history for ports.
type Tracker struct {
	mu     sync.Mutex
	states map[string]*State
	clock  func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		states: make(map[string]*State),
		clock:  time.Now,
	}
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 8)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

// Update records whether a port is currently up or down.
func (t *Tracker) Update(port int, proto string, up bool) *State {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	k := key(port, proto)
	s, exists := t.states[k]
	if !exists {
		s = &State{Port: port, Proto: proto, FirstSeen: now}
		t.states[k] = s
	}
	if exists && s.Up != up {
		s.Flaps++
	}
	s.Up = up
	s.LastSeen = now
	return s
}

// Get returns the state for a port, and whether it exists.
func (t *Tracker) Get(port int, proto string) (*State, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.states[key(port, proto)]
	return s, ok
}

// All returns a snapshot of all tracked states.
func (t *Tracker) All() []*State {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]*State, 0, len(t.states))
	for _, s := range t.states {
		copy := *s
		out = append(out, &copy)
	}
	return out
}

// Reset clears all tracked state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.states = make(map[string]*State)
}
