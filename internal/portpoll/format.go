package portpoll

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes a human-readable poll result table to stdout.
func Print(results []Result) { Fprint(os.Stdout, results) }

// Fprint writes a human-readable poll result table to w.
func Fprint(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no poll results")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ADDR\tPROTO\tSTATUS\tLATENCY")

	for _, r := range results {
		status := "open"
		if !r.Open {
			status = "closed"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Addr, r.Proto, status, r.Latency.Round(1000))
	}
	tw.Flush()
}
