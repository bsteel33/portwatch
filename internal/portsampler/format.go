package portsampler

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Print writes a summary of samples to stdout.
func Print(samples []Sample) { Fprint(os.Stdout, samples) }

// Fprint writes a summary of samples to w.
func Fprint(w io.Writer, samples []Sample) {
	if len(samples) == 0 {
		fmt.Fprintln(w, "no samples recorded")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tOPEN PORTS")
	for _, s := range samples {
		fmt.Fprintf(tw, "%s\t%d\n", s.At.Format(time.RFC3339), s.Count)
	}
	tw.Flush()
}

// PrintLast writes only the most recent sample to stdout.
func PrintLast(s Sample) { FprintLast(os.Stdout, s) }

// FprintLast writes only the most recent sample to w.
func FprintLast(w io.Writer, s Sample) {
	fmt.Fprintf(w, "last sample at %s: %d open port(s)\n", s.At.Format(time.RFC3339), s.Count)
}
