package portrank

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Print writes a ranked entry table to stdout.
func Print(entries []Entry) { Fprint(os.Stdout, entries) }

// Fprint writes a ranked entry table to w.
func Fprint(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no ports to rank")
		return
	}
	fmt.Fprintf(w, "%-8s %-6s %-10s %s\n", "PORT", "PROTO", "RISK", "SCORE")
	fmt.Fprintln(w, strings.Repeat("-", 34))
	for _, e := range entries {
		fmt.Fprintf(w, "%-8d %-6s %-10s %d\n",
			e.Port.Port, e.Port.Proto, e.Level.String(), e.Score)
	}
}
