// Package notify provides notification delivery for portwatch.
//
// Supported channels:
//
//   - stdout: writes a human-readable diff summary to standard output (default).
//   - exec:   invokes an external command, passing the diff summary as the
//             first argument. Useful for integrating with custom alerting
//             scripts or third-party tools.
//
// Basic usage:
//
//	cfg := notify.DefaultConfig()
//	n := notify.New(cfg)
//	if err := n.Notify(diff); err != nil {
//	    log.Println("notify error:", err)
//	}
package notify
