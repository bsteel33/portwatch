// Package porttrend analyses the directional trend of open port counts over
// a configurable sliding time window.
//
// Usage:
//
//	tracker := porttrend.New(porttrend.DefaultConfig())
//	tracker.Record(len(ports))
//	// … later …
//	result := tracker.Analyze()
//	porttrend.Print(result)
//
// Directions:
//   - "up"     – port count is growing beyond the configured threshold
//   - "down"   – port count is shrinking beyond the configured threshold
//   - "stable" – change is within the threshold
package porttrend
