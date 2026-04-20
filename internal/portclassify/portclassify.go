// Package portclassify assigns a risk classification to open ports
// based on well-known vulnerability patterns and service type.
package portclassify

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/scanner"
)

// Class represents a risk classification level.
type Class int

const (
	ClassUnknown  Class = iota
	ClassSafe            // expected, low-risk services
	ClassMonitor         // worth watching but not inherently risky
	ClassSuspicious      // unusual or commonly exploited
	ClassDangerous       // high-risk, should rarely be exposed
)

func (c Class) String() string {
	switch c {
	case ClassSafe:
		return "safe"
	case ClassMonitor:
		return "monitor"
	case ClassSuspicious:
		return "suspicious"
	case ClassDangerous:
		return "dangerous"
	default:
		return "unknown"
	}
}

// Result holds the classification for a single port.
type Result struct {
	Port   scanner.Port
	Class  Class
	Reason string
}

// Classifier classifies ports by risk level.
type Classifier struct {
	cfg Config
}

// New returns a new Classifier with the given config.
func New(cfg Config) *Classifier {
	return &Classifier{cfg: cfg}
}

// Classify evaluates each port and returns a Result slice.
func (c *Classifier) Classify(ports []scanner.Port) []Result {
	results := make([]Result, 0, len(ports))
	for _, p := range ports {
		results = append(results, c.classifyOne(p))
	}
	return results
}

func (c *Classifier) classifyOne(p scanner.Port) Result {
	if cls, reason, ok := matchRules(c.cfg.Rules, p); ok {
		return Result{Port: p, Class: cls, Reason: reason}
	}
	return Result{Port: p, Class: ClassUnknown, Reason: "no matching rule"}
}

func matchRules(rules []Rule, p scanner.Port) (Class, string, bool) {
	for _, r := range rules {
		if r.Proto != "" && r.Proto != p.Proto {
			continue
		}
		if r.Port != 0 && r.Port != p.Port {
			continue
		}
		return r.Class, r.Reason, true
	}
	return ClassUnknown, "", false
}

// Fprint writes a human-readable classification table to w.
func Fprint(w io.Writer, results []Result) {
	for _, r := range results {
		fmt.Fprintf(w, "%-6d %-5s %-12s %s\n", r.Port.Port, r.Port.Proto, r.Class.String(), r.Reason)
	}
}

// Print writes the classification table to stdout.
func Print(results []Result) { Fprint(os.Stdout, results) }
