// Package trend provides a sliding-window port count tracker for portwatch.
//
// It records successive open-port counts and exposes a Delta method to detect
// whether the number of open ports is growing or shrinking over a configurable
// time window. Useful for alerting on sustained port-count changes rather than
// individual scan-to-scan fluctuations.
package trend
