// Package portguard provides a port allowlist/denylist guard that enforces
// a static policy on observed ports, flagging any that violate the policy.
package portguard

import (
	"fmt"
	"sync"

	"github.com/example/portwatch/internal/scanner"
)

// Action represents the guard decision for a port.
type Action string

const (
	ActionAllow  Action = "allow"
	ActionDeny   Action = "deny"
	ActionUnknown Action = "unknown"
)

// Violation describes a port that violated the guard policy.
type Violation struct {
	Port    scanner.Port
	Action  Action
	Reason  string
}

// Guard enforces an allowlist or denylist policy over a set of ports.
type Guard struct {
	mu       sync.RWMutex
	allowlist map[string]struct{}
	denylist  map[string]struct{}
	default_  Action
}

// New creates a Guard with the given config.
func New(cfg Config) *Guard {
	g := &Guard{
		allowlist: make(map[string]struct{}),
		denylist:  make(map[string]struct{}),
		default_:  cfg.Default,
	}
	for _, k := range cfg.Allowlist {
		g.allowlist[k] = struct{}{}
	}
	for _, k := range cfg.Denylist {
		g.denylist[k] = struct{}{}
	}
	return g
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Port, p.Proto)
}

// Evaluate returns the Action and optional Violation for a single port.
func (g *Guard) Evaluate(p scanner.Port) (Action, *Violation) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	k := portKey(p)

	if _, ok := g.denylist[k]; ok {
		v := &Violation{Port: p, Action: ActionDeny, Reason: "port is explicitly denied"}
		return ActionDeny, v
	}
	if _, ok := g.allowlist[k]; ok {
		return ActionAllow, nil
	}
	if g.default_ == ActionDeny {
		v := &Violation{Port: p, Action: ActionDeny, Reason: "port not in allowlist (default deny)"}
		return ActionDeny, v
	}
	return ActionAllow, nil
}

// Check evaluates all ports and returns any violations.
func (g *Guard) Check(ports []scanner.Port) []Violation {
	var out []Violation
	for _, p := range ports {
		_, v := g.Evaluate(p)
		if v != nil {
			out = append(out, *v)
		}
	}
	return out
}
