// Package suppress provides alert suppression for portwatch.
//
// When a port change triggers an alert, the Suppressor can be used to
// prevent repeated notifications for the same port key within a
// configurable TTL window. Once the TTL expires the key is eligible
// for alerting again.
//
// Usage:
//
//	cfg := suppress.DefaultConfig()
//	s := suppress.New(cfg)
//	if !s.IsSuppressed(key) {
//		sendAlert()
//		s.Suppress(key)
//	}
package suppress
