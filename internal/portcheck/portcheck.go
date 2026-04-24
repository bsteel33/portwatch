// Package portcheck evaluates a set of ports against user-defined health
// conditions and reports which ports fail the check.
package portcheck

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/scanner"
)

// Condition describes a single health check rule.
type Condition struct {
	Port  int
	Proto string // "tcp" or "udp"
	MustBeOpen bool
}

// Result holds the outcome of a single condition evaluation.
type Result struct {
	Condition Condition
	Passed    bool
	Reason    string
}

// Checker evaluates port health conditions.
type Checker struct {
	conditions []Condition
}

// New creates a Checker with the given conditions.
func New(conditions []Condition) *Checker {
	return &Checker{conditions: conditions}
}

// Evaluate checks each condition against the provided open ports and returns
// one Result per condition.
func (c *Checker) Evaluate(ports []scanner.Port) []Result {
	index := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		index[key(p.Port, p.Proto)] = struct{}{}
	}

	results := make([]Result, 0, len(c.conditions))
	for _, cond := range c.conditions {
		_, open := index[key(cond.Port, cond.Proto)]
		passed := open == cond.MustBeOpen
		reason := ""
		if !passed {
			if cond.MustBeOpen {
				reason = fmt.Sprintf("port %d/%s expected open but is closed", cond.Port, cond.Proto)
			} else {
				reason = fmt.Sprintf("port %d/%s expected closed but is open", cond.Port, cond.Proto)
			}
		}
		results = append(results, Result{Condition: cond, Passed: passed, Reason: reason})
	}
	return results
}

// AnyFailed returns true if at least one result did not pass.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if !r.Passed {
			return true
		}
	}
	return false
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Print writes a human-readable summary of results to stdout.
func Print(results []Result) { Fprint(os.Stdout, results) }

// Fprint writes a human-readable summary of results to w.
func Fprint(w io.Writer, results []Result) {
	for _, r := range results {
		status := "OK"
		if !r.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(w, "[%s] %d/%s: %s\n", status, r.Condition.Port, r.Condition.Proto, func() string {
			if r.Reason != "" {
				return r.Reason
			}
			return "check passed"
		}())
	}
}
