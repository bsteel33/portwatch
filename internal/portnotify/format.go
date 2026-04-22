package portnotify

import (
	"fmt"
	"io"
	"os"
)

// Print writes all events to stdout.
func Print(events []Event) {
	Fprint(os.Stdout, events)
}

// Fprint writes all events to w.
func Fprint(w io.Writer, events []Event) {
	if len(events) == 0 {
		fmt.Fprintln(w, "portnotify: no matching events")
		return
	}
	fmt.Fprintf(w, "portnotify: %d event(s)\n", len(events))
	for _, e := range events {
		fmt.Fprintf(w, "  rule=%-20s port=%5d proto=%-4s service=%s\n",
			e.Rule.Label, e.Port.Port, e.Port.Proto, serviceName(e.Port.Service))
	}
}

func serviceName(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
