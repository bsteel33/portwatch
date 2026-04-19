package portmap

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Print writes all mappings to stdout.
func Print(m *Map) { Fprint(os.Stdout, m) }

// Fprint writes all mappings to w.
func Fprint(w io.Writer, m *Map) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.entries))
	for k := range m.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "PORT MAP")
	fmt.Fprintln(w, "--------")
	if len(keys) == 0 {
		fmt.Fprintln(w, "  (empty)")
		return
	}
	for _, k := range keys {
		e := m.entries[k]
		fmt.Fprintf(w, "  %s/%d  =>  %s\n", e.Proto, e.Port, e.Name)
	}
}
