// Package portshield provides a simple allow/block shield that gates
// port visibility based on a configured set of trusted port numbers.
// Ports not explicitly shielded are passed through unchanged.
package portshield

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Shield holds the set of protected ports and their action.
type Shield struct {
	mu      sync.RWMutex
	rules   map[string]Action
	default_ Action
}

// Action describes what the shield does when a port matches.
type Action int

const (
	Allow Action = iota
	Block
)

func (a Action) String() string {
	if a == Allow {
		return "allow"
	}
	return "block"
}

// New creates a Shield with the given default action.
func New(defaultAction Action) *Shield {
	return &Shield{
		rules:    make(map[string]Action),
		default_: defaultAction,
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Add registers an action for a specific port/proto combination.
func (s *Shield) Add(port int, proto string, action Action) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[portKey(port, proto)] = action
}

// Evaluate returns the Action that applies to the given port.
func (s *Shield) Evaluate(port scanner.Port) Action {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if a, ok := s.rules[portKey(port.Port, port.Proto)]; ok {
		return a
	}
	return s.default_
}

// Filter returns only the ports that are Allowed by the shield.
func (s *Shield) Filter(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		if s.Evaluate(p) == Allow {
			out = append(out, p)
		}
	}
	return out
}
