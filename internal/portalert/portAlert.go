// Package portAlert evaluates open ports against a set of alert rules
// and emits Alert values for any port that matches.
package portAlert

import (
	"fmt"
	"net"
)

// Severity classifies how urgent an alert is.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Rule describes a single alerting condition.
type Rule struct {
	Port     int
	Proto    string // "tcp" or "udp"; empty matches both
	Severity Severity
	Message  string
}

// Alert is emitted when a scanned port matches a rule.
type Alert struct {
	Addr     string
	Proto    string
	Port     int
	Severity Severity
	Message  string
}

// Evaluator holds a set of rules and evaluates ports against them.
type Evaluator struct {
	rules []Rule
}

// New returns an Evaluator loaded with the provided rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate checks each (proto, port) pair in addrs against the rule set.
// addrs is expected in "proto:host:port" or "proto::port" form; use
// net.JoinHostPort for the address part.
type Port struct {
	Proto string
	Port  int
}

// Check evaluates a slice of Port values and returns matching alerts.
func (e *Evaluator) Check(ports []Port) []Alert {
	var alerts []Alert
	for _, p := range ports {
		for _, r := range e.rules {
			if r.Port != p.Port {
				continue
			}
			if r.Proto != "" && r.Proto != p.Proto {
				continue
			}
			alerts = append(alerts, Alert{
				Addr:     net.JoinHostPort("localhost", fmt.Sprintf("%d", p.Port)),
				Proto:    p.Proto,
				Port:     p.Port,
				Severity: r.Severity,
				Message:  r.Message,
			})
		}
	}
	return alerts
}
