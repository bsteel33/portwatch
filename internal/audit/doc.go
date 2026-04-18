// Package audit provides persistent audit trail logging for portwatch.
//
// Each scan that produces a change (added or removed ports) is recorded as a
// JSON entry in an append-only log file. The log can be inspected with the
// built-in formatter or read programmatically via Load.
//
// Usage:
//
//	logger := audit.New("/var/log/portwatch-audit.log")
//	logger.Record("alert", diff)
package audit
