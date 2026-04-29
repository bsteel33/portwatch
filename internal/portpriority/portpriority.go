// Package portpriority assigns priority levels to ports based on
// configurable rules, allowing operators to triage alerts by importance.
package portpriority

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents a priority level for a port.
type Level int

const (
	Low      Level = 1
	Medium   Level = 2
	High     Level = 3
	Critical Level = 4
)

func (l Level) String() string {
	switch l {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	case Critical:
		return "critical"
	default:
		return "unknown"
	}
}

// Rule maps a port/proto pair to a priority level.
type Rule struct {
	Port  int
	Proto string
	Level Level
}

// Prioritizer assigns priority levels to open ports.
type Prioritizer struct {
	mu      sync.RWMutex
	rules   []Rule
	default_ Level
}

// New creates a Prioritizer with the given rules and a fallback default level.
func New(rules []Rule, defaultLevel Level) *Prioritizer {
	return &Prioritizer{
		rules:   rules,
		default_: defaultLevel,
	}
}

// Assign returns the priority level for the given port.
func (p *Prioritizer) Assign(port scanner.Port) Level {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.rules {
		if r.Port == port.Port && (r.Proto == "" || r.Proto == port.Proto) {
			return r.Level
		}
	}
	return p.default_
}

// AssignAll returns a map of port key to Level for a slice of ports.
func (p *Prioritizer) AssignAll(ports []scanner.Port) map[string]Level {
	out := make(map[string]Level, len(ports))
	for _, port := range ports {
		k := fmt.Sprintf("%d/%s", port.Port, port.Proto)
		out[k] = p.Assign(port)
	}
	return out
}
