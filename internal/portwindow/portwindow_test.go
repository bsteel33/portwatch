package portwindow

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func fixedTime(hour, minute int) func() time.Time {
	return func() time.Time {
		return time.Date(2024, 1, 1, hour, minute, 0, 0, time.UTC)
	}
}

func mustParseWindow(t *testing.T, s string) Window {
	t.Helper()
	w, err := ParseWindow(s)
	if err != nil {
		t.Fatalf("ParseWindow(%q): %v", s, err)
	}
	return w
}

func TestActive_NoWindows(t *testing.T) {
	pw := New(Config{})
	pw.now = fixedTime(3, 0)
	if !pw.Active() {
		t.Error("expected Active=true when no windows configured")
	}
}

func TestActive_WithinWindow(t *testing.T) {
	pw := New(Config{Windows: []Window{mustParseWindow(t, "09:00-17:00")}})
	pw.now = fixedTime(12, 30)
	if !pw.Active() {
		t.Error("expected Active=true at 12:30 within 09:00-17:00")
	}
}

func TestActive_OutsideWindow(t *testing.T) {
	pw := New(Config{Windows: []Window{mustParseWindow(t, "09:00-17:00")}})
	pw.now = fixedTime(18, 0)
	if pw.Active() {
		t.Error("expected Active=false at 18:00 outside 09:00-17:00")
	}
}

func TestActive_MidnightWraparound(t *testing.T) {
	pw := New(Config{Windows: []Window{mustParseWindow(t, "22:00-06:00")}})
	pw.now = fixedTime(23, 30)
	if !pw.Active() {
		t.Error("expected Active=true at 23:30 within 22:00-06:00")
	}
	pw.now = fixedTime(3, 0)
	if !pw.Active() {
		t.Error("expected Active=true at 03:00 within 22:00-06:00")
	}
	pw.now = fixedTime(10, 0)
	if pw.Active() {
		t.Error("expected Active=false at 10:00 outside 22:00-06:00")
	}
}

func TestFilter_ActiveReturnsAll(t *testing.T) {
	ports := []scanner.Port{{Port: 80, Proto: "tcp"}, {Port: 443, Proto: "tcp"}}
	pw := New(Config{Windows: []Window{mustParseWindow(t, "00:00-23:59")}})
	pw.now = fixedTime(12, 0)
	got := pw.Filter(ports)
	if len(got) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestFilter_InactiveReturnsNil(t *testing.T) {
	ports := []scanner.Port{{Port: 80, Proto: "tcp"}}
	pw := New(Config{Windows: []Window{mustParseWindow(t, "09:00-10:00")}})
	pw.now = fixedTime(22, 0)
	got := pw.Filter(ports)
	if got != nil {
		t.Errorf("expected nil ports outside window, got %v", got)
	}
}

func TestParseWindow_Invalid(t *testing.T) {
	_, err := ParseWindow("not-a-window")
	if err == nil {
		t.Error("expected error for invalid window string")
	}
}
