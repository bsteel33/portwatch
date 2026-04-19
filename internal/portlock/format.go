package portlock

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes the locked ports table to stdout.
func Print(l *Locker) {
	Fprint(os.Stdout, l)
}

// Fprint writes the locked ports table to w.
func Fprint(w io.Writer, l *Locker) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.locks) == 0 {
		fmt.Fprintln(w, "no locked ports")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTO\tNOTE")
	for _, lk := range l.locks {
		fmt.Fprintf(tw, "%d\t%s\t%s\n", lk.Port, lk.Proto, lk.Note)
	}
	tw.Flush()
}
