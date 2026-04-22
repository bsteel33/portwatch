// Package portnotify delivers per-port notifications when a port's state
// crosses a configured threshold or matches a watch rule.
package portnotify

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Event describes a single port notification.
type Event struct {
	Port    scanner.Port
	Rule    Rule
	Message string
}

// Notifier evaluates open ports against watch rules and emits events.
type Notifier struct {
	cfg    Config
	rules  []Rule
	out    io.Writer
}

// New creates a Notifier with the given config.
func New(cfg Config) *Notifier {
	return &Notifier{
		cfg:   cfg,
		rules: cfg.Rules,
		out:   os.Stdout,
	}
}

// Check evaluates ports against all rules and returns matching events.
func (n *Notifier) Check(ports []scanner.Port) []Event {
	var events []Event
	for _, p := range ports {
		for _, r := range n.rules {
			if matchesRule(p, r) {
				events = append(events, Event{
					Port:    p,
					Rule:    r,
					Message: buildMsg(p, r),
				})
			}
		}
	}
	return events
}

// Notify writes events to the configured output.
func (n *Notifier) Notify(events []Event) {
	for _, e := range events {
		fmt.Fprintln(n.out, e.Message)
	}
}

func matchesRule(p scanner.Port, r Rule) bool {
	if r.Port != 0 && r.Port != p.Port {
		return false
	}
	if r.Proto != "" && !strings.EqualFold(r.Proto, p.Proto) {
		return false
	}
	return true
}

func buildMsg(p scanner.Port, r Rule) string {
	svc := p.Service
	if svc == "" {
		svc = "unknown"
	}
	return fmt.Sprintf("[portnotify] rule=%q port=%d proto=%s service=%s",
		r.Label, p.Port, p.Proto, svc)
}
