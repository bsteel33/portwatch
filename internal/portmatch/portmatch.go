// Package portmatch provides pattern-based port matching using glob-style rules.
package portmatch

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Rule represents a single match rule (e.g. "80/tcp", "*/udp", "8080").
type Rule struct {
	Port  string // "*" or numeric
	Proto string // "*", "tcp", or "udp"
}

// Matcher holds a set of rules and matches ports against them.
type Matcher struct {
	rules []Rule
}

// New creates a Matcher from a slice of raw rule strings.
func New(raw []string) (*Matcher, error) {
	rules := make([]Rule, 0, len(raw))
	for _, s := range raw {
		r, err := ParseRule(s)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return &Matcher{rules: rules}, nil
}

// ParseRule parses a rule string into a Rule.
// Accepted formats: "80", "80/tcp", "*/udp", "*".
func ParseRule(s string) (Rule, error) {
	parts := strings.SplitN(s, "/", 2)
	port := parts[0]
	proto := "*"
	if len(parts) == 2 {
		proto = strings.ToLower(parts[1])
		if proto != "tcp" && proto != "udp" && proto != "*" {
			return Rule{}, fmt.Errorf("portmatch: invalid protocol %q in rule %q", proto, s)
		}
	}
	if port != "*" {
		if _, err := strconv.Atoi(port); err != nil {
			return Rule{}, fmt.Errorf("portmatch: invalid port %q in rule %q", port, s)
		}
	}
	return Rule{Port: port, Proto: proto}, nil
}

// Match returns true if p matches any rule in the Matcher.
func (m *Matcher) Match(p scanner.Port) bool {
	for _, r := range m.rules {
		if matchRule(r, p) {
			return true
		}
	}
	return false
}

// Filter returns only the ports that match at least one rule.
func (m *Matcher) Filter(ports []scanner.Port) []scanner.Port {
	if len(m.rules) == 0 {
		return ports
	}
	out := make([]scanner.Port, 0)
	for _, p := range ports {
		if m.Match(p) {
			out = append(out, p)
		}
	}
	return out
}

func matchRule(r Rule, p scanner.Port) bool {
	portMatch := r.Port == "*" || r.Port == strconv.Itoa(p.Port)
	protoMatch := r.Proto == "*" || r.Proto == strings.ToLower(p.Proto)
	return portMatch && protoMatch
}
