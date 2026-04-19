// Package portreport aggregates per-port metadata into a unified report entry.
package portreport

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds aggregated metadata for a single port.
type Entry struct {
	Port    int
	Proto   string
	Service string
	State   string
	Tags    []string
	Score   int
}

// Report is an ordered collection of port entries.
type Report struct {
	entries []Entry
}

// New builds a Report from a slice of open ports.
func New(ports []scanner.Port, svc func(int, string) string) *Report {
	r := &Report{}
	for _, p := range ports {
		name := ""
		if svc != nil {
			name = svc(p.Port, p.Proto)
		}
		r.entries = append(r.entries, Entry{
			Port:    p.Port,
			Proto:   p.Proto,
			Service: name,
			State:   "open",
		})
	}
	sort.Slice(r.entries, func(i, j int) bool {
		if r.entries[i].Port != r.entries[j].Port {
			return r.entries[i].Port < r.entries[j].Port
		}
		return r.entries[i].Proto < r.entries[j].Proto
	})
	return r
}

// Entries returns all report entries.
func (r *Report) Entries() []Entry { return r.entries }

// Len returns the number of entries.
func (r *Report) Len() int { return len(r.entries) }

// Fprint writes a human-readable table to w.
func Fprint(w io.Writer, r *Report) {
	if r.Len() == 0 {
		fmt.Fprintln(w, "no open ports")
		return
	}
	fmt.Fprintf(w, "%-8s %-6s %-20s %s\n", "PORT", "PROTO", "SERVICE", "STATE")
	for _, e := range r.entries {
		svc := e.Service
		if svc == "" {
			svc = "unknown"
		}
		fmt.Fprintf(w, "%-8d %-6s %-20s %s\n", e.Port, e.Proto, svc, e.State)
	}
}

// Print writes the report to stdout.
func Print(r *Report) { Fprint(os.Stdout, r) }
