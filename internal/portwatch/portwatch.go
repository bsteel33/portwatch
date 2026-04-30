// Package portwatch provides per-port watch rules that trigger callbacks
// when specific ports appear, disappear, or change state.
package portwatch

import (
	"fmt"
	"sync"

	"github.com/example/portwatch/internal/scanner"
)

// Event describes what happened to a watched port.
type Event struct {
	Port   scanner.Port
	Rule   Rule
	Reason string
}

// Watcher evaluates a set of Rules against a port list and fires callbacks.
type Watcher struct {
	mu    sync.Mutex
	cfg   Config
	rules []Rule
}

// New returns a Watcher using cfg.
func New(cfg Config) *Watcher {
	return &Watcher{cfg: cfg}
}

// AddRule appends a rule to the watcher.
func (w *Watcher) AddRule(r Rule) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rules = append(w.rules, r)
}

// Evaluate checks ports against all rules and returns matching events.
func (w *Watcher) Evaluate(ports []scanner.Port) []Event {
	w.mu.Lock()
	defer w.mu.Unlock()

	var events []Event
	for _, p := range ports {
		for _, r := range w.rules {
			if matchesRule(p, r) {
				events = append(events, Event{
					Port:   p,
					Rule:   r,
					Reason: fmt.Sprintf("port %d/%s matched rule %q", p.Port, p.Proto, r.Name),
				})
			}
		}
	}
	return events
}

func matchesRule(p scanner.Port, r Rule) bool {
	if r.Port != 0 && r.Port != p.Port {
		return false
	}
	if r.Proto != "" && r.Proto != p.Proto {
		return false
	}
	return true
}
