package portprobe

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Fprint writes a human-readable summary of all probe results to w.
func Fprint(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no probe results")
		return
	}
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Port != sorted[j].Port {
			return sorted[i].Port < sorted[j].Port
		}
		return sorted[i].Proto < sorted[j].Proto
	})
	fmt.Fprintf(w, "%-8s %-6s %-12s %s\n", "PORT", "PROTO", "REACHABLE", "LATENCY")
	for _, r := range sorted {
		reach := "yes"
		if !r.Reachable {
			reach = "no"
		}
		fmt.Fprintf(w, "%-8d %-6s %-12s %s\n", r.Port, r.Proto, reach, r.Latency.Round(1000))
	}
}

// Print writes probe results to stdout.
func Print(results []Result) {
	Fprint(os.Stdout, results)
}
