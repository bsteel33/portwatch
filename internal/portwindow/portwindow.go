// Package portwindow provides time-window-based filtering of port scan results,
// allowing ports to be included or excluded based on configurable time-of-day windows.
package portwindow

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Window defines an active time range during which ports should be monitored.
type Window struct {
	Start time.Time // time-of-day start (only hour/minute used)
	End   time.Time // time-of-day end (only hour/minute used)
}

// Config holds configuration for the port window filter.
type Config struct {
	Windows  []Window
	Location *time.Location
}

// PortWindow filters scan results based on active time windows.
type PortWindow struct {
	cfg Config
	now func() time.Time
}

// New creates a new PortWindow with the given config.
func New(cfg Config) *PortWindow {
	if cfg.Location == nil {
		cfg.Location = time.Local
	}
	return &PortWindow{cfg: cfg, now: time.Now}
}

// Active reports whether the current time falls within any configured window.
func (pw *PortWindow) Active() bool {
	if len(pw.cfg.Windows) == 0 {
		return true
	}
	now := pw.now().In(pw.cfg.Location)
	h, m, _ := now.Clock()
	current := h*60 + m
	for _, w := range pw.cfg.Windows {
		sh, sm, _ := w.Start.Clock()
		eh, em, _ := w.End.Clock()
		start := sh*60 + sm
		end := eh*60 + em
		if start <= end {
			if current >= start && current < end {
				return true
			}
		} else {
			// wraps midnight
			if current >= start || current < end {
				return true
			}
		}
	}
	return false
}

// Filter returns ports only if the current time is within an active window.
// If no windows are configured, all ports are returned unchanged.
func (pw *PortWindow) Filter(ports []scanner.Port) []scanner.Port {
	if pw.Active() {
		return ports
	}
	return nil
}

// ParseWindow parses a window string in "HH:MM-HH:MM" format.
func ParseWindow(s string) (Window, error) {
	var sh, sm, eh, em int
	_, err := fmt.Sscanf(s, "%d:%d-%d:%d", &sh, &sm, &eh, &em)
	if err != nil {
		return Window{}, fmt.Errorf("portwindow: invalid window %q: expected HH:MM-HH:MM", s)
	}
	base := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	return Window{
		Start: base.Add(time.Duration(sh)*time.Hour + time.Duration(sm)*time.Minute),
		End:   base.Add(time.Duration(eh)*time.Hour + time.Duration(em)*time.Minute),
	}, nil
}
