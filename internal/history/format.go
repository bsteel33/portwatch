package history

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a human-readable summary of all history entries to w.
func Print(h *History, w io.Writer) {
	if len(h.Entries) == 0 {
		fmt.Fprintln(w, "no history recorded")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tADDED\tREMOVED")
	for _, e := range h.Entries {
		fmt.Fprintf(tw, "%s\t+%d\t-%d\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			len(e.Added),
			len(e.Removed),
		)
	}
	tw.Flush()
}

// PrintEntry writes detailed info for a single Entry to w.
func PrintEntry(e *Entry, w io.Writer) {
	fmt.Fprintf(w, "Time: %s\n", e.Timestamp.Format("2006-01-02 15:04:05"))
	if len(e.Added) > 0 {
		fmt.Fprintln(w, "Added:")
		for _, p := range e.Added {
			fmt.Fprintf(w, "  + %d/%s (%s)\n", p.Port, p.Proto, p.Service)
		}
	}
	if len(e.Removed) > 0 {
		fmt.Fprintln(w, "Removed:")
		for _, p := range e.Removed {
			fmt.Fprintf(w, "  - %d/%s (%s)\n", p.Port, p.Proto, p.Service)
		}
	}
}
