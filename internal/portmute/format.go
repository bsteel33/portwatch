package portmute

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// Print writes all active mute entries to stdout.
func Print(m *Muter) {
	Fprint(os.Stdout, m)
}

// Fprint writes all active mute entries to w.
func Fprint(w io.Writer, m *Muter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()

	type row struct {
		k      string
		entry  Entry
		remain time.Duration
	}

	var rows []row
	for k, e := range m.entries {
		if now.After(e.Until) {
			continue
		}
		rows = append(rows, row{k: k, entry: e, remain: e.Until.Sub(now).Truncate(time.Second)})
	}

	if len(rows) == 0 {
		fmt.Fprintln(w, "no muted ports")
		return
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].k < rows[j].k })

	fmt.Fprintf(w, "%-20s %-12s %s\n", "PORT", "REMAINING", "REASON")
	for _, r := range rows {
		reason := r.entry.Reason
		if reason == "" {
			reason = "-"
		}
		fmt.Fprintf(w, "%-20s %-12s %s\n", r.k, r.remain, reason)
	}
}
