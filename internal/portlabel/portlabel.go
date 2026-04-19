// Package portlabel assigns human-readable labels to ports based on
// configurable rules and well-known service mappings.
package portlabel

import (
	"fmt"
	"sync"
)

// Label represents a named tag for a port.
type Label struct {
	Port  int
	Proto string
	Name  string
}

// Labeler maps ports to labels.
type Labeler struct {
	mu     sync.RWMutex
	labels map[string]string
}

// New returns a Labeler seeded with the provided labels.
func New(initial []Label) *Labeler {
	l := &Labeler{labels: make(map[string]string)}
	for _, lb := range initial {
		l.labels[key(lb.Port, lb.Proto)] = lb.Name
	}
	return l
}

// Set adds or updates a label for a port/proto pair.
func (l *Labeler) Set(port int, proto, name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.labels[key(port, proto)] = name
}

// Get returns the label for a port/proto pair and whether it was found.
func (l *Labeler) Get(port int, proto string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	v, ok := l.labels[key(port, proto)]
	return v, ok
}

// Remove deletes the label for a port/proto pair.
func (l *Labeler) Remove(port int, proto string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.labels, key(port, proto))
}

// All returns a copy of all current labels.
func (l *Labeler) All() []Label {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Label, 0, len(l.labels))
	for k, name := range l.labels {
		var port int
		var proto string
		fmt.Sscanf(k, "%d/%s", &port, &proto)
		out = append(out, Label{Port: port, Proto: proto, Name: name})
	}
	return out
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
