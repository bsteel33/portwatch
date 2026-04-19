// Package filter provides port filtering based on include and exclude rules.
//
// Rules can match by port number, protocol, or both. When include rules are
// present, only ports matching at least one include rule are kept. Exclude
// rules are always applied and take precedence over include rules.
//
// A Rule with neither Port nor Protocol set matches all ports (wildcard).
//
// Example usage:
//
//	cfg := filter.Config{
//		Include: []filter.Rule{{Protocol: "tcp"}},
//		Exclude: []filter.Rule{{Port: 22}},
//	}
//	f := filter.New(cfg)
//	filtered := f.Apply(scannedPorts)
//
// To exclude a specific port regardless of protocol:
//
//	cfg := filter.Config{
//		Exclude: []filter.Rule{{Port: 22}, {Port: 23}},
//	}
package filter
