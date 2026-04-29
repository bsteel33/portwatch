// Package portrelay provides fan-out forwarding of port scan snapshots to
// multiple downstream destinations.
//
// Destinations implement the Destination interface and can be registered at
// runtime.  Included implementations:
//
//   - WriterDestination — writes JSON to any io.Writer (stdout, file, …)
//   - HTTPDestination   — POSTs JSON to an HTTP endpoint
//
// Example:
//
//	relay := portrelay.New(os.Stderr)
//	relay.Register(portrelay.NewWriterDestination("stdout", os.Stdout))
//	relay.Forward(snap)
package portrelay
