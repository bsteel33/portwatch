package portflag

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Fprint writes all flagged ports and their flags to w.
func Fprint(w io.Writer, f *Flagger) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.flags) == 0 {
		fmt.Fprintln(w, "no port flags set")
		return
	}

	keys := make([]string, 0, len(f.flags))
	for k := range f.flags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "%-20s  %s\n", "PORT/PROTO", "FLAGS")
	fmt.Fprintf(w, "%-20s  %s\n", "----------", "-----")
	for _, k := range keys {
		names := make([]string, 0, len(f.flags[k]))
		for name := range f.flags[k] {
			names = append(names, name)
		}
		sort.Strings(names)
		fmt.Fprintf(w, "%-20s  %v\n", k, names)
	}
}

// Print writes all flagged ports to stdout.
func Print(f *Flagger) {
	Fprint(os.Stdout, f)
}
