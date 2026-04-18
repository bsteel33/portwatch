package snapshot

import (
	"fmt"
	"io"
	"strings"
)

// PrintDiff writes a human-readable diff summary to w.
func PrintDiff(w io.Writer, d *Diff) {
	if !d.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return
	}
	if len(d.Opened) > 0 {
		fmt.Fprintln(w, "Opened ports:")
		for _, p := range d.Opened {
			fmt.Fprintf(w, "  + %d/%s (%s)\n", p.Port, p.Proto, p.Service)
		}
	}
	if len(d.Closed) > 0 {
		fmt.Fprintln(w, "Closed ports:")
		for _, p := range d.Closed {
			fmt.Fprintf(w, "  - %d/%s (%s)\n", p.Port, p.Proto, p.Service)
		}
	}
}

// PrintSnapshot writes a human-readable snapshot summary to w.
func PrintSnapshot(w io.Writer, s *Snapshot) {
	fmt.Fprintf(w, "Snapshot taken at %s\n", s.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))
	if len(s.Ports) == 0 {
		fmt.Fprintln(w, "No open ports found.")
		return
	}
	for _, p := range s.Ports {
		fmt.Fprintf(w, "  %d/%s\t%s\n", p.Port, p.Proto, p.Service)
	}
}
