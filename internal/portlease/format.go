package portlease

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// Fprint writes active leases to w in a human-readable table.
func Fprint(w io.Writer, leases map[string]Lease, now time.Time) {
	if len(leases) == 0 {
		fmt.Fprintln(w, "no active leases")
		return
	}

	keys := make([]string, 0, len(leases))
	for k := range leases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT/PROTO\tOWNER\tEXPIRES IN")
	for _, k := range keys {
		lease := leases[k]
		remaining := lease.ExpiresAt.Sub(now).Round(time.Second)
		fmt.Fprintf(tw, "%s\t%s\t%s\n", k, lease.Owner, remaining)
	}
	_ = tw.Flush()
}

// Print writes active leases to stdout.
func Print(leases map[string]Lease, now time.Time) {
	Fprint(os.Stdout, leases, now)
}
