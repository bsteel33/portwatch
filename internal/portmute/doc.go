// Package portmute implements per-port alert muting for portwatch.
//
// A port can be muted for a configurable duration so that transient or
// expected open-port events do not trigger noisy alerts. Mute entries
// expire automatically; they can also be lifted early with Unmute.
//
// Typical usage:
//
//	m := portmute.New()
//	m.Mute(8080, "tcp", 10*time.Minute, "maintenance window")
//	if m.IsMuted(8080, "tcp") {
//		// skip alert
//	}
package portmute
