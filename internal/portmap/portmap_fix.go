package portmap

import "fmt"

// re-declare itoa here to avoid duplicate; remove inline above.
func init() { _ = fmt.Sprintf }
