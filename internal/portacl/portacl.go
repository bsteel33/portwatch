// Package portacl implements an access control list for ports,
// allowing or denying ports based on configured rules.
package portacl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Action represents an ACL decision.
type Action string

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

// Rule is a single ACL entry.
type Rule struct {
	Port  int
	Proto string
	Action Action
}

// ACL evaluates ports against an ordered list of rules.
type ACL struct {
	rules []Rule
	defaultAction Action
}

// New creates an ACL with the given rules and a default action applied
// when no rule matches.
func New(rules []Rule, defaultAction Action) *ACL {
	return &ACL{rules: rules, defaultAction: defaultAction}
}

// Evaluate returns the Action for the given port.
func (a *ACL) Evaluate(p scanner.Port) Action {
	for _, r := range a.rules {
		if r.Port != 0 && r.Port != p.Port {
			continue
		}
		if r.Proto != "" && !strings.EqualFold(r.Proto, p.Proto) {
			continue
		}
		return r.Action
	}
	return a.defaultAction
}

// Filter returns only ports that are allowed by the ACL.
func (a *ACL) Filter(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		if a.Evaluate(p) == Allow {
			out = append(out, p)
		}
	}
	return out
}

// ParseRule parses a rule string of the form "allow:80", "deny:23/tcp".
func ParseRule(s string) (Rule, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("portacl: invalid rule %q", s)
	}
	action := Action(strings.ToLower(parts[0]))
	if action != Allow && action != Deny {
		return Rule{}, fmt.Errorf("portacl: unknown action %q", parts[0])
	}
	portProto := parts[1]
	var port int
	var proto string
	if idx := strings.Index(portProto, "/"); idx >= 0 {
		proto = strings.ToLower(portProto[idx+1:])
		portProto = portProto[:idx]
	}
	if portProto != "*" && portProto != "" {
		v, err := strconv.Atoi(portProto)
		if err != nil {
			return Rule{}, fmt.Errorf("portacl: invalid port %q", portProto)
		}
		port = v
	}
	return Rule{Port: port, Proto: proto, Action: action}, nil
}
