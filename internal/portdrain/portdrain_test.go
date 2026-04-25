package portdrain

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

var (
	fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	p80      = scanner.Port{Port: 80, Proto: "tcp"}
	p443     = scanner.Port{Port: 443, Proto: "tcp"}
)

func newTestDrainer(now time.Time) *Drainer {
	d := New()
	d.now = func() time.Time { return now }
	return d
}

func TestMark_And_IsDraining(t *testing.T) {
	d := newTestDrainer(fixedNow)
	d.Mark(p80, 5*time.Minute)

	if !d.IsDraining(p80) {
		t.Fatal("expected port 80 to be draining")
	}
	if d.IsDraining(p443) {
		t.Fatal("expected port 443 not to be draining")
	}
}

func TestIsDraining_Expired(t *testing.T) {
	d := newTestDrainer(fixedNow)
	d.Mark(p80, 1*time.Minute)

	// advance clock past deadline
	d.now = func() time.Time { return fixedNow.Add(2 * time.Minute) }

	if d.IsDraining(p80) {
		t.Fatal("expected drain entry to be expired")
	}
}

func TestOverdue_ReturnsExpiredOpenPorts(t *testing.T) {
	d := newTestDrainer(fixedNow)
	d.Mark(p80, 1*time.Minute)
	d.Mark(p443, 10*time.Minute)

	d.now = func() time.Time { return fixedNow.Add(2 * time.Minute) }

	overdue := d.Overdue([]scanner.Port{p80, p443})
	if len(overdue) != 1 {
		t.Fatalf("expected 1 overdue port, got %d", len(overdue))
	}
	if overdue[0].Port.Port != 80 {
		t.Errorf("expected port 80 to be overdue, got %d", overdue[0].Port.Port)
	}
}

func TestOverdue_ClosedPortNotReturned(t *testing.T) {
	d := newTestDrainer(fixedNow)
	d.Mark(p80, 1*time.Minute)

	d.now = func() time.Time { return fixedNow.Add(2 * time.Minute) }

	// p80 is NOT in the open set — it already closed
	overdue := d.Overdue([]scanner.Port{p443})
	if len(overdue) != 0 {
		t.Fatalf("expected no overdue ports, got %d", len(overdue))
	}
}

func TestEvict_RemovesEntry(t *testing.T) {
	d := newTestDrainer(fixedNow)
	d.Mark(p80, 5*time.Minute)
	d.Evict(p80)

	if d.IsDraining(p80) {
		t.Fatal("expected port 80 to be removed after eviction")
	}
	if len(d.All()) != 0 {
		t.Fatal("expected drain list to be empty after eviction")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	d := newTestDrainer(fixedNow)
	d.Mark(p80, 5*time.Minute)
	d.Mark(p443, 10*time.Minute)

	all := d.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}
