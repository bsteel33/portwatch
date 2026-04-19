package portannot

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Print writes all annotations to stdout.
func Print(a *Annotator) {
	Fprint(os.Stdout, a)
}

// Fprint writes all annotations to w.
func Fprint(w io.Writer, a *Annotator) {
	list := a.All()
	if len(list) == 0 {
		fmt.Fprintln(w, "no annotations")
		return
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Port != list[j].Port {
			return list[i].Port < list[j].Port
		}
		return list[i].Proto < list[j].Proto
	})
	for _, ann := range list {
		fmt.Fprintf(w, "%-6d %-5s %s\n", ann.Port, ann.Proto, ann.Note)
	}
}
