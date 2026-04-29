// Package portversion tracks the banner or version string reported by services
// on open ports. On each scan the caller supplies the observed version string;
// the tracker persists the last-seen value and surfaces a Change whenever the
// string differs from the previously recorded one.
//
// Typical usage:
//
//	tracker, _ := portversion.New(cfg.Path)
//	for _, p := range openPorts {
//		if ch := tracker.Update(p.Port, p.Proto, p.Banner); ch != nil {
//			portversion.PrintChange(ch)
//		}
//	}
//	tracker.Save()
package portversion
