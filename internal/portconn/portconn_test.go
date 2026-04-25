package portconn

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	tr := New(0)
	tr.Set(80, "tcp", 5)

	e, ok := tr.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Count != 5 {
		t.Errorf("expected count 5, got %d", e.Count)
	}
	if e.Port != 80 || e.Proto != "tcp" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := New(0)
	_, ok := tr.Get(443, "tcp")
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestExceeded_BelowThreshold(t *testing.T) {
	tr := New(10)
	tr.Set(80, "tcp", 5)
	tr.Set(443, "tcp", 3)

	if got := tr.Exceeded(); len(got) != 0 {
		t.Errorf("expected no exceeded entries, got %d", len(got))
	}
}

func TestExceeded_AboveThreshold(t *testing.T) {
	tr := New(10)
	tr.Set(80, "tcp", 15)
	tr.Set(443, "tcp", 3)
	tr.Set(22, "tcp", 20)

	got := tr.Exceeded()
	if len(got) != 2 {
		t.Fatalf("expected 2 exceeded entries, got %d", len(got))
	}
	// sorted by port: 22, 80
	if got[0].Port != 22 {
		t.Errorf("expected port 22 first, got %d", got[0].Port)
	}
	if got[1].Port != 80 {
		t.Errorf("expected port 80 second, got %d", got[1].Port)
	}
}

func TestExceeded_DisabledThreshold(t *testing.T) {
	tr := New(0)
	tr.Set(80, "tcp", 999)

	if got := tr.Exceeded(); got != nil {
		t.Errorf("expected nil when threshold disabled, got %v", got)
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	tr := New(5)
	tr.Set(80, "tcp", 10)
	tr.Reset()

	_, ok := tr.Get(80, "tcp")
	if ok {
		t.Fatal("expected empty tracker after reset")
	}
	if got := tr.Exceeded(); len(got) != 0 {
		t.Errorf("expected no exceeded entries after reset")
	}
}
