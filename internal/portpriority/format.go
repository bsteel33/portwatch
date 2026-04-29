package portpriority

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Print writes a prioritized port table to stdout.
func Print(ports []scanner.Port, p *Prioritizer) {
	Fprint(os.Stdout, ports, p)
}

// Fprint writes a prioritized port table to w.
func Fprint(w io.Writer, ports []scanner.Port, p *Prioritizer) {
	type row struct {
		port  scanner.Port
		level Level
	}
	rows := make([]row, len(ports))
	for i, port := range ports {
		rows[i] = row{port: port, level: p.Assign(port)}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].level != rows[j].level {
			return rows[i].level > rows[j].level
		}
		return rows[i].port.Port < rows[j].port.Port
	})
	fmt.Fprintf(w, "%-8s %-6s %-10s %s\n", "PRIORITY", "PORT", "PROTO", "SERVICE")
	for _, r := range rows {
		fmt.Fprintf(w, "%-8s %-6d %-10s %s\n",
			r.level.String(),
			r.port.Port,
			r.port.Proto,
			r.port.Service,
		)
	}
}
