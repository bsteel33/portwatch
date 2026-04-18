package audit

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes a human-readable audit log to stdout.
func Print(path string) error {
	return Fprint(os.Stdout, path)
}

// Fprint writes a human-readable audit log to w.
func Fprint(w io.Writer, path string) error {
	entries, err := Load(path)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "no audit entries found")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tEVENT\tADDED\tREMOVED")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%d\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Event,
			len(e.Added),
			len(e.Removed),
		)
	}
	return tw.Flush()
}

// PrintEntry writes detail for a single entry.
func PrintEntry(w io.Writer, e Entry) {
	fmt.Fprintf(w, "[%s] event=%s added=%d removed=%d\n",
		e.Timestamp.Format("2006-01-02T15:04:05Z"),
		e.Event, len(e.Added), len(e.Removed))
	for _, p := range e.Added {
		fmt.Fprintf(w, "  + %d/%s\n", p.Port, p.Proto)
	}
	for _, p := range e.Removed {
		fmt.Fprintf(w, "  - %d/%s\n", p.Port, p.Proto)
	}
}
