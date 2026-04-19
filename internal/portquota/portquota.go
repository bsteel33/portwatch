// Package portquota enforces per-protocol and total open port count limits.
package portquota

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Violation describes a quota breach.
type Violation struct {
	Proto string
	Limit int
	Actual int
}

func (v Violation) Error() string {
	return fmt.Sprintf("quota exceeded: proto=%s limit=%d actual=%d", v.Proto, v.Limit, v.Actual)
}

// Quota holds limits for port counts.
type Quota struct {
	mu     sync.Mutex
	cfg    Config
	totals map[string]int // proto -> count
	total  int
}

// New creates a Quota enforcer with the given config.
func New(cfg Config) *Quota {
	return &Quota{cfg: cfg, totals: make(map[string]int)}
}

// Check evaluates ports against configured limits and returns all violations.
func (q *Quota) Check(ports []scanner.Port) []Violation {
	q.mu.Lock()
	defer q.mu.Unlock()

	counts := make(map[string]int)
	for _, p := range ports {
		counts[p.Proto]++
	}

	q.totals = counts
	q.total = len(ports)

	var violations []Violation

	if q.cfg.TotalLimit > 0 && q.total > q.cfg.TotalLimit {
		violations = append(violations, Violation{Proto: "any", Limit: q.cfg.TotalLimit, Actual: q.total})
	}

	for proto, limit := range q.cfg.ProtoLimits {
		if actual, ok := counts[proto]; ok && actual > limit {
			violations = append(violations, Violation{Proto: proto, Limit: limit, Actual: actual})
		}
	}

	return violations
}

// Totals returns the last computed per-protocol counts.
func (q *Quota) Totals() map[string]int {
	q.mu.Lock()
	defer q.mu.Unlock()
	copy := make(map[string]int, len(q.totals))
	for k, v := range q.totals {
		copy[k] = v
	}
	return copy
}
