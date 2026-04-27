package portjournal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const timeLayout = "2006-01-02 15:04:05"

// Print writes all journal entries to stdout.
func Print(j *Journal) { Fprint(os.Stdout, j) }

// Fprint writes all journal entries to w.
func Fprint(w io.Writer, j *Journal) {
	entries := j.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "(no journal entries)")
		return
	}
	fmt.Fprintf(w, "%-20s  %-6s  %-5s  %-12s  %s\n",
		"TIME", "KIND", "PORT", "PROTO", "SERVICE")
	fmt.Fprintln(w, strings.Repeat("-", 62))
	for _, e := range entries {
		printEntry(w, e)
	}
}

// PrintLast writes the n most recent entries to stdout.
func PrintLast(j *Journal, n int) { FprintLast(os.Stdout, j, n) }

// FprintLast writes the n most recent entries to w.
func FprintLast(w io.Writer, j *Journal, n int) {
	entries := j.Last(n)
	if len(entries) == 0 {
		fmt.Fprintln(w, "(no journal entries)")
		return
	}
	for _, e := range entries {
		printEntry(w, e)
	}
}

func printEntry(w io.Writer, e Entry) {
	svc := e.Service
	if svc == "" {
		svc = "-"
	}
	note := ""
	if e.Note != "" {
		note = " (" + e.Note + ")"
	}
	fmt.Fprintf(w, "%-20s  %-6s  %-5d  %-12s  %s%s\n",
		e.Time.Format(timeLayout),
		string(e.Kind),
		e.Port,
		e.Proto,
		svc,
		note,
	)
}
