package baseline

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes a human-readable baseline summary to stdout.
func Print(b *Baseline) {
	Fprint(os.Stdout, b)
}

// Fprint writes a human-readable baseline summary to w.
func Fprint(w io.Writer, b *Baseline) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Baseline captured at: %s\n", b.CapturedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(tw, "PORT\tPROTO\tSERVICE\n")
	for _, p := range b.Ports {
		fmt.Fprintf(tw, "%d\t%s\t%s\n", p.Port, p.Proto, p.Service)
	}
	tw.Flush()
}

// PrintDeviation writes deviation details to stdout.
func PrintDeviation(d Deviation) {
	FprintDeviation(os.Stdout, d)
}

// FprintDeviation writes deviation details to w.
func FprintDeviation(w io.Writer, d Deviation) {
	if !d.HasChanges() {
		fmt.Fprintln(w, "No deviations from baseline.")
		return
	}
	for _, p := range d.Added {
		fmt.Fprintf(w, "[+] %d/%s %s (not in baseline)\n", p.Port, p.Proto, p.Service)
	}
	for _, p := range d.Removed {
		fmt.Fprintf(w, "[-] %d/%s %s (missing from current)\n", p.Port, p.Proto, p.Service)
	}
}
