package portbatch

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Fprint writes a human-readable summary of a Batch to w.
func Fprint(w io.Writer, b Batch) {
	if len(b.Ports) == 0 {
		fmt.Fprintln(w, "batch: (empty)")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "batch collected at %s — %d port(s)\n",
		b.Collected.Format("15:04:05"), len(b.Ports))
	fmt.Fprintln(tw, "PORT\tPROTO\tSERVICE")
	for _, p := range b.Ports {
		fmt.Fprintf(tw, "%d\t%s\t%s\n", p.Port, p.Proto, p.Service)
	}
	_ = tw.Flush()
}

// Print writes a human-readable summary of a Batch to stdout.
func Print(b Batch) {
	Fprint(os.Stdout, b)
}
