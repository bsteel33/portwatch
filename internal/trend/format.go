package trend

import (
	"fmt"
	"io"
	"os"
)

// Print writes a human-readable trend summary to stdout.
func Print(t *Trend) {
	Fprint(os.Stdout, t)
}

// Fprint writes a human-readable trend summary to w.
func Fprint(w io.Writer, t *Trend) {
	pts := t.Points()
	if len(pts) == 0 {
		fmt.Fprintln(w, "trend: no data")
		return
	}
	first := pts[0]
	last := pts[len(pts)-1]
	delta := last.Count - first.Count
	sign := "+"
	if delta < 0 {
		sign = ""
	}
	fmt.Fprintf(w, "trend: %d points over %s | first=%d last=%d delta=%s%d\n",
		len(pts),
		last.At.Sub(first.At).Round(1e6),
		first.Count,
		last.Count,
		sign,
		delta,
	)
}
