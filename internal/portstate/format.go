package portstate

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Print writes all tracked states to stdout.
func Print(t *Tracker) {
	Fprint(os.Stdout, t)
}

// Fprint writes all tracked states to w.
func Fprint(w io.Writer, t *Tracker) {
	states := t.All()
	sort.Slice(states, func(i, j int) bool {
		if states[i].Port != states[j].Port {
			return states[i].Port < states[j].Port
		}
		return states[i].Proto < states[j].Proto
	})
	for _, s := range states {
		status := "down"
		if s.Up {
			status = "up"
		}
		fmt.Fprintf(w, "%-6d %-5s %-4s flaps=%-3d first=%-20s last=%s\n",
			s.Port, s.Proto, status, s.Flaps,
			s.FirstSeen.Format("2006-01-02 15:04:05"),
			s.LastSeen.Format("2006-01-02 15:04:05"),
		)
	}
}
