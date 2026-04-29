package portage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portage.json")
}

func newTracker(t *testing.T, now func() time.Time) *Tracker {
	t.Helper()
	tr, err := New(tempPath(t), now)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

var (
	t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 = t0.Add(5 * time.Minute)
)

func TestUpdate_RecordsFirstSeen(t *testing.T) {
	tr := newTracker(t, func() time.Time { return t0 })
	ports := []scanner.Port{{Port: 80, Proto: "tcp"}}
	if err := tr.Update(ports); err != nil {
		t.Fatalf("Update: %v", err)
	}
	e := tr.Get(ports[0])
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if !e.FirstSeen.Equal(t0) {
		t.Errorf("FirstSeen = %v, want %v", e.FirstSeen, t0)
	}
}

func TestUpdate_DoesNotOverwriteFirstSeen(t *testing.T) {
	clock := t0
	tr := newTracker(t, func() time.Time { return clock })
	ports := []scanner.Port{{Port: 443, Proto: "tcp"}}
	_ = tr.Update(ports)
	clock = t1
	_ = tr.Update(ports)
	e := tr.Get(ports[0])
	if !e.FirstSeen.Equal(t0) {
		t.Errorf("FirstSeen overwritten: got %v, want %v", e.FirstSeen, t0)
	}
	if !e.LastSeen.Equal(t1) {
		t.Errorf("LastSeen = %v, want %v", e.LastSeen, t1)
	}
}

func TestUpdate_EvictsClosedPort(t *testing.T) {
	tr := newTracker(t, func() time.Time { return t0 })
	p := scanner.Port{Port: 22, Proto: "tcp"}
	_ = tr.Update([]scanner.Port{p})
	_ = tr.Update([]scanner.Port{})
	if tr.Get(p) != nil {
		t.Error("expected eviction, entry still present")
	}
}

func TestAge_ReturnsElapsed(t *testing.T) {
	clock := t0
	tr := newTracker(t, func() time.Time { return clock })
	p := scanner.Port{Port: 8080, Proto: "tcp"}
	_ = tr.Update([]scanner.Port{p})
	clock = t0.Add(10 * time.Minute)
	age, ok := tr.Age(p)
	if !ok {
		t.Fatal("expected age, got not-found")
	}
	if age != 10*time.Minute {
		t.Errorf("Age = %v, want 10m", age)
	}
}

func TestNew_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	tr, err := New(path, nil)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestPersistence(t *testing.T) {
	path := tempPath(t)
	tr, _ := New(path, func() time.Time { return t0 })
	p := scanner.Port{Port: 3306, Proto: "tcp"}
	_ = tr.Update([]scanner.Port{p})

	tr2, err := New(path, func() time.Time { return t1 })
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e := tr2.Get(p)
	if e == nil {
		t.Fatal("entry missing after reload")
	}
	if !e.FirstSeen.Equal(t0) {
		t.Errorf("FirstSeen after reload = %v, want %v", e.FirstSeen, t0)
	}
}

func TestNew_CorruptFile(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := New(path, nil)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
