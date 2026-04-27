package portpool

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Fprint writes a human-readable summary of the pool to w.
func Fprint(w io.Writer, p *Pool) {
	ports := p.All()
	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Port != ports[j].Port {
			return ports[i].Port < ports[j].Port
		}
		return ports[i].Proto < ports[j].Proto
	})

	cap := "unlimited"
	if p.Capacity() > 0 {
		cap = fmt.Sprintf("%d", p.Capacity())
	}
	fmt.Fprintf(w, "Pool: %s  entries: %d/%s\n", p.Name(), p.Len(), cap)
	for _, port := range ports {
		svc := port.Service
		if svc == "" {
			svc = "unknown"
		}
		fmt.Fprintf(w, "  %5d/%-4s  %s\n", port.Port, port.Proto, svc)
	}
}

// Print writes a human-readable summary of the pool to stdout.
func Print(p *Pool) {
	Fprint(os.Stdout, p)
}
