// Package throttle provides a simple per-key rate limiter used by portwatch
// to suppress duplicate alerts when the same port change is detected across
// multiple consecutive scans.
//
// Usage:
//
//	cfg := throttle.DefaultConfig()
//	th := throttle.New(cfg)
//	if th.Allow("tcp:8080:added") {
//		// send alert
//	}
package throttle
