package portexpiry

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// Print writes expired entries to stdout.
func Print(entries []Entry, max time.Duration) {
	Fprint(os.Stdout, entries, max)
}

// Fprint writes expired entries to w.
func Fprint(w io.Writer, entries []Entry, max time.Duration) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no expired ports")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].FirstSeen.Before(entries[j].FirstSeen)
	})
	fmt.Fprintf(w, "%-8s %-6s %-30s %s\n", "PORT", "PROTO", "FIRST SEEN", "AGE")
	for _, e := range entries {
		age := time.Since(e.FirstSeen).Round(time.Second)
		fmt.Fprintf(w, "%-8d %-6s %-30s %s  [EXPIRED > %s]\n",
			e.Port, e.Proto, e.FirstSeen.Format(time.RFC3339), age, max)
	}
}
