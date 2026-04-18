// Package fingerprint generates a stable hash for a set of open ports,
// allowing quick equality checks between snapshots.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a hex-encoded SHA-256 digest of a port set.
type Fingerprint string

// Compute returns a deterministic Fingerprint for the given ports.
// Ports are sorted before hashing so order does not affect the result.
func Compute(ports []scanner.Port) Fingerprint {
	sorted := make([]scanner.Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Proto != sorted[j].Proto {
			return sorted[i].Proto < sorted[j].Proto
		}
		return sorted[i].Number < sorted[j].Number
	})

	h := sha256.New()
	for _, p := range sorted {
		fmt.Fprintf(h, "%s:%d\n", p.Proto, p.Number)
	}
	return Fingerprint(hex.EncodeToString(h.Sum(nil)))
}

// Equal returns true when two Fingerprints are identical.
func Equal(a, b Fingerprint) bool {
	return a == b
}
