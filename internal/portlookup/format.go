package portlookup

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes a formatted table of lookup results to stdout.
func Print(results []Result) {
	Fprint(os.Stdout, results)
}

// Fprint writes a formatted table of lookup results to w.
func Fprint(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no results")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTO\tSERVICE\tRESOLVED")
	for _, r := range results {
		resolvedStr := "no"
		if r.Found {
			resolvedStr = "yes"
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", r.Port, r.Proto, r.Service, resolvedStr)
	}
	_ = tw.Flush()
}
