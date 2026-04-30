// Package portwatch implements per-port watch rules for portwatch.
//
// A Watcher holds a set of Rules, each identifying a port number and optional
// protocol. Calling Evaluate against a live port list returns an Event for
// every port that satisfies a rule, allowing callers to react to specific
// ports appearing on the host.
//
// Example:
//
//	w := portwatch.New(portwatch.DefaultConfig())
//	w.AddRule(portwatch.Rule{Name: "ssh", Port: 22, Proto: "tcp"})
//	events := w.Evaluate(ports)
package portwatch
