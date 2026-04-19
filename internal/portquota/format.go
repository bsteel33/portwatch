package portquota

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes violations to stdout.
func Print(violations []Violation) { Fprint(os.Stdout, violations) }

// Fprint writes violations to w.
func Fprint(w io.Writer, violations []Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "portquota: no violations")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROTO\tLIMIT\tACTUAL")
	for _, v := range violations {
		fmt.Fprintf(tw, "%s\t%d\t%d\n", v.Proto, v.Limit, v.Actual)
	}
	tw.Flush()
}
