package portversion

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Fprint writes all tracked version entries to w in a human-readable table.
func Fprint(w io.Writer, t *Tracker) {
	t.mu.Lock()
	defer t.mu.Unlock()

	keys := make([]string, 0, len(t.entries))
	for k := range t.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "%-20s  %s\n", "PORT/PROTO", "VERSION")
	fmt.Fprintf(w, "%-20s  %s\n", "----------", "-------")
	for _, k := range keys {
		e := t.entries[k]
		v := e.Version
		if v == "" {
			v = "(unknown)"
		}
		fmt.Fprintf(w, "%-20s  %s\n", fmt.Sprintf("%d/%s", e.Port, e.Proto), v)
	}
}

// Print writes all tracked version entries to stdout.
func Print(t *Tracker) {
	Fprint(os.Stdout, t)
}

// FprintChange writes a single Change to w.
func FprintChange(w io.Writer, c *Change) {
	fmt.Fprintf(w, "version change on %d/%s: %q -> %q\n",
		c.Port, c.Proto, c.OldVersion, c.NewVersion)
}

// PrintChange writes a single Change to stdout.
func PrintChange(c *Change) {
	FprintChange(os.Stdout, c)
}
