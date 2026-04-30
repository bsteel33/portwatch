package portwatch

import (
	"fmt"
	"io"
	"os"
)

// Print writes events to stdout.
func Print(events []Event) {
	Fprint(os.Stdout, events)
}

// Fprint writes events to w.
func Fprint(w io.Writer, events []Event) {
	if len(events) == 0 {
		fmt.Fprintln(w, "portwatch: no events")
		return
	}
	fmt.Fprintf(w, "portwatch: %d event(s)\n", len(events))
	for _, e := range events {
		fmt.Fprintf(w, "  [%s] port=%d proto=%s reason=%q\n",
			e.Rule.Name, e.Port.Port, e.Port.Proto, e.Reason)
	}
}
