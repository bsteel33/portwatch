// Package portschedule provides time-window-based scheduling for port scan
// operations. It allows scans to be enabled or suppressed based on named
// schedules (e.g. "business-hours", "maintenance-window"), each defined as
// a set of weekday + time-of-day ranges.
package portschedule

import (
	"fmt"
	"strings"
	"time"
)

// Weekday maps a lowercase weekday name to time.Weekday.
var weekdayNames = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

// Window represents a single time window within a day.
type Window struct {
	Weekday time.Weekday
	Start   time.Duration // offset from midnight
	End     time.Duration // offset from midnight
}

// Schedule holds a named collection of time windows.
type Schedule struct {
	Name    string
	Windows []Window
}

// Scheduler manages multiple named schedules and evaluates whether a given
// moment falls within any active schedule.
type Scheduler struct {
	schedules map[string]*Schedule
	clock     func() time.Time
}

// New returns a Scheduler with the given clock function. Pass time.Now for
// production use; inject a fixed clock in tests.
func New(clock func() time.Time) *Scheduler {
	return &Scheduler{
		schedules: make(map[string]*Schedule),
		clock:     clock,
	}
}

// Add registers a named schedule. Duplicate names are overwritten.
func (s *Scheduler) Add(sch *Schedule) {
	s.schedules[sch.Name] = sch
}

// Active returns true if the current time falls within any window of the
// named schedule. Returns false and no error when the schedule is not found.
func (s *Scheduler) Active(name string) (bool, error) {
	sch, ok := s.schedules[name]
	if !ok {
		return false, fmt.Errorf("portschedule: schedule %q not found", name)
	}
	now := s.clock()
	return sch.contains(now), nil
}

// ActiveAny returns true if the current time falls within any window of any
// registered schedule.
func (s *Scheduler) ActiveAny() bool {
	now := s.clock()
	for _, sch := range s.schedules {
		if sch.contains(now) {
			return true
		}
	}
	return false
}

// Names returns the names of all registered schedules in insertion order.
func (s *Scheduler) Names() []string {
	out := make([]string, 0, len(s.schedules))
	for k := range s.schedules {
		out = append(out, k)
	}
	return out
}

// contains reports whether t falls within any of the schedule's windows.
func (sch *Schedule) contains(t time.Time) bool {
	offset := todayOffset(t)
	for _, w := range sch.Windows {
		if t.Weekday() == w.Weekday && offset >= w.Start && offset < w.End {
			return true
		}
	}
	return false
}

// todayOffset returns the duration since midnight for the given time.
func todayOffset(t time.Time) time.Duration {
	h, m, sec := t.Clock()
	return time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(sec)*time.Second
}

// ParseWindow parses a window string of the form "monday 09:00-17:00".
// The day name is case-insensitive.
func ParseWindow(s string) (Window, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return Window{}, fmt.Errorf("portschedule: invalid window %q: expected \"day HH:MM-HH:MM\"", s)
	}
	wd, ok := weekdayNames[strings.ToLower(parts[0])]
	if !ok {
		return Window{}, fmt.Errorf("portschedule: unknown weekday %q", parts[0])
	}
	rangeParts := strings.SplitN(parts[1], "-", 2)
	if len(rangeParts) != 2 {
		return Window{}, fmt.Errorf("portschedule: invalid time range %q", parts[1])
	}
	start, err := parseHHMM(rangeParts[0])
	if err != nil {
		return Window{}, fmt.Errorf("portschedule: %w", err)
	}
	end, err := parseHHMM(rangeParts[1])
	if err != nil {
		return Window{}, fmt.Errorf("portschedule: %w", err)
	}
	if end <= start {
		return Window{}, fmt.Errorf("portschedule: end time must be after start time in %q", s)
	}
	return Window{Weekday: wd, Start: start, End: end}, nil
}

// parseHHMM converts "HH:MM" into a duration from midnight.
func parseHHMM(s string) (time.Duration, error) {
	var h, m int
	_, err := fmt.Sscanf(s, "%d:%d", &h, &m)
	if err != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid time %q: expected HH:MM", s)
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}
