package portmemo

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Fprint writes all memos from s to w in a human-readable table.
func Fprint(w io.Writer, s *Store) {
	memos := s.All()
	if len(memos) == 0 {
		fmt.Fprintln(w, "no memos stored")
		return
	}
	sort.Slice(memos, func(i, j int) bool {
		if memos[i].Port != memos[j].Port {
			return memos[i].Port < memos[j].Port
		}
		return memos[i].Proto < memos[j].Proto
	})
	fmt.Fprintf(w, "%-8s %-6s %-20s %s\n", "PORT", "PROTO", "KEY", "VALUE")
	fmt.Fprintf(w, "%-8s %-6s %-20s %s\n", "----", "-----", "-------------------", "-----")
	for _, m := range memos {
		keys := make([]string, 0, len(m.Notes))
		for k := range m.Notes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "%-8d %-6s %-20s %s\n", m.Port, m.Proto, k, m.Notes[k])
		}
	}
}

// Print writes all memos to stdout.
func Print(s *Store) {
	Fprint(os.Stdout, s)
}
