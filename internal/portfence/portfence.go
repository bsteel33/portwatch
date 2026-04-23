// Package portfence enforces port access boundaries by blocking or allowing
// scan results based on configurable IP/CIDR range rules.
package portfence

import (
	"fmt"
	"net"
	"strings"
)

// Action represents the fence policy for a rule.
type Action string

const (
	ActionAllow Action = "allow"
	ActionBlock  Action = "block"
)

// Rule defines a single fence rule pairing a CIDR with an action.
type Rule struct {
	CIDR   *net.IPNet
	Action Action
}

// Fence evaluates IP addresses against an ordered list of rules.
type Fence struct {
	rules      []Rule
	defaultAct Action
}

// New creates a Fence with the given rules and a fallback default action.
func New(rules []Rule, defaultAction Action) *Fence {
	return &Fence{rules: rules, defaultAct: defaultAction}
}

// Evaluate returns the Action that applies to the given IP address.
// Rules are evaluated in order; the first match wins. If no rule matches,
// the default action is returned.
func (f *Fence) Evaluate(ip net.IP) Action {
	for _, r := range f.rules {
		if r.CIDR.Contains(ip) {
			return r.Action
		}
	}
	return f.defaultAct
}

// Allowed reports whether the given IP is permitted under the fence policy.
func (f *Fence) Allowed(ip net.IP) bool {
	return f.Evaluate(ip) == ActionAllow
}

// ParseRule parses a rule string of the form "<cidr>:<action>",
// e.g. "192.168.0.0/16:allow" or "10.0.0.0/8:block".
func ParseRule(s string) (Rule, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("portfence: invalid rule %q: expected <cidr>:<action>", s)
	}
	_, cidr, err := net.ParseCIDR(strings.TrimSpace(parts[0]))
	if err != nil {
		return Rule{}, fmt.Errorf("portfence: invalid CIDR in rule %q: %w", s, err)
	}
	act := Action(strings.TrimSpace(strings.ToLower(parts[1])))
	if act != ActionAllow && act != ActionBlock {
		return Rule{}, fmt.Errorf("portfence: unknown action %q in rule %q", act, s)
	}
	return Rule{CIDR: cidr, Action: act}, nil
}
