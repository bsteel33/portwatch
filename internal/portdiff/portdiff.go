// Package portdiff provides a human-readable summary of port scan differences.
package portdiff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry represents a single line in a diff summary.
type Entry struct {
	Port  int
	Proto string
	Name  string
	Kind  string // "added" | "removed"
}

// Summary holds the result of diffing two snapshots.
type Summary struct {
	Added   []Entry
	Removed []Entry
}

// Build computes a Summary from two snapshots.
func Build(prev, curr *snapshot.Snapshot) Summary {
	prevIdx := index(prev)
	currIdx := index(curr)

	var s Summary
	for k, e := range currIdx {
		if _, ok := prevIdx[k]; !ok {
			s.Added = append(s.Added, Entry{Port: e.Port, Proto: e.Proto, Name: e.Name, Kind: "added"})
		}
	}
	for k, e := range prevIdx {
		if _, ok := currIdx[k]; !ok {
			s.Removed = append(s.Removed, Entry{Port: e.Port, Proto: e.Proto, Name: e.Name, Kind: "removed"})
		}
	}
	return s
}

// HasChanges reports whether the summary contains any differences.
func (s Summary) HasChanges() bool {
	return len(s.Added) > 0 || len(s.Removed) > 0
}

// Fprint writes the summary to w.
func Fprint(w io.Writer, s Summary) {
	if !s.HasChanges() {
		fmt.Fprintln(w, "no changes detected")
		return
	}
	for _, e := range s.Added {
		fmt.Fprintf(w, "+ %d/%s %s\n", e.Port, e.Proto, label(e.Name))
	}
	for _, e := range s.Removed {
		fmt.Fprintf(w, "- %d/%s %s\n", e.Port, e.Proto, label(e.Name))
	}
}

// Print writes the summary to stdout.
func Print(s Summary) { Fprint(os.Stdout, s) }

func label(name string) string {
	if strings.TrimSpace(name) == "" {
		return "(unknown)"
	}
	return name
}

type portEntry struct {
	Port  int
	Proto string
	Name  string
}

func index(snap *snapshot.Snapshot) map[string]portEntry {
	m := make(map[string]portEntry)
	if snap == nil {
		return m
	}
	for _, p := range snap.Ports {
		k := fmt.Sprintf("%d/%s", p.Port, p.Proto)
		m[k] = portEntry{Port: p.Port, Proto: p.Proto, Name: p.Name}
	}
	return m
}
