package portstate

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newTestTracker() *Tracker {
	t := New()
	t.clock = func() time.Time { return fixedNow }
	return t
}

func TestUpdate_NewPort(t *testing.T) {
	tr := newTestTracker()
	s := tr.Update(80, "tcp", true)
	if !s.Up {
		t.Fatal("expected port to be up")
	}
	if s.Flaps != 0 {
		t.Fatalf("expected 0 flaps, got %d", s.Flaps)
	}
	if s.FirstSeen != fixedNow {
		t.Fatalf("unexpected FirstSeen: %v", s.FirstSeen)
	}
}

func TestUpdate_Flap(t *testing.T) {
	tr := newTestTracker()
	tr.Update(443, "tcp", true)
	tr.Update(443, "tcp", false)
	s := tr.Update(443, "tcp", true)
	if s.Flaps != 2 {
		t.Fatalf("expected 2 flaps, got %d", s.Flaps)
	}
}

func TestUpdate_NoFlapSameState(t *testing.T) {
	tr := newTestTracker()
	tr.Update(22, "tcp", true)
	s := tr.Update(22, "tcp", true)
	if s.Flaps != 0 {
		t.Fatalf("expected 0 flaps, got %d", s.Flaps)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := newTestTracker()
	_, ok := tr.Get(9999, "tcp")
	if ok {
		t.Fatal("expected missing port to return false")
	}
}

func TestGet_Found(t *testing.T) {
	tr := newTestTracker()
	tr.Update(8080, "tcp", true)
	s, ok := tr.Get(8080, "tcp")
	if !ok {
		t.Fatal("expected port to be found")
	}
	if s.Port != 8080 {
		t.Fatalf("expected port 8080, got %d", s.Port)
	}
}

func TestReset(t *testing.T) {
	tr := newTestTracker()
	tr.Update(80, "tcp", true)
	tr.Reset()
	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker after reset")
	}
}

func TestAll_ReturnsCopies(t *testing.T) {
	tr := newTestTracker()
	tr.Update(80, "tcp", true)
	all := tr.All()
	all[0].Flaps = 99
	s, _ := tr.Get(80, "tcp")
	if s.Flaps == 99 {
		t.Fatal("All() should return copies, not references")
	}
}
