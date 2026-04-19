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

var fixed = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newTracker(t *testing.T) *Tracker {
	t.Helper()
	tr, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	tr.now = func() time.Time { return fixed }
	return tr
}

func TestUpdate_RecordsFirstSeen(t *testing.T) {
	tr := newTracker(t)
	ports := []scanner.Port{{Port: 80, Proto: "tcp"}}
	tr.Update(ports)
	e, ok := tr.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if !e.FirstSeen.Equal(fixed) {
		t.Errorf("FirstSeen = %v, want %v", e.FirstSeen, fixed)
	}
}

func TestUpdate_DoesNotOverwriteFirstSeen(t *testing.T) {
	tr := newTracker(t)
	ports := []scanner.Port{{Port: 443, Proto: "tcp"}}
	tr.Update(ports)
	later := fixed.Add(time.Hour)
	tr.now = func() time.Time { return later }
	tr.Update(ports)
	e, _ := tr.Get(443, "tcp")
	if !e.FirstSeen.Equal(fixed) {
		t.Errorf("FirstSeen overwritten: got %v, want %v", e.FirstSeen, fixed)
	}
}

func TestUpdate_EvictsClosedPort(t *testing.T) {
	tr := newTracker(t)
	tr.Update([]scanner.Port{{Port: 22, Proto: "tcp"}})
	tr.Update([]scanner.Port{})
	if _, ok := tr.Get(22, "tcp"); ok {
		t.Error("expected port 22 to be evicted")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	tr, _ := New(path)
	tr.now = func() time.Time { return fixed }
	tr.Update([]scanner.Port{{Port: 8080, Proto: "tcp"}})
	if err := tr.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	tr2, err := New(path)
	if err != nil {
		t.Fatalf("New (reload): %v", err)
	}
	e, ok := tr2.Get(8080, "tcp")
	if !ok {
		t.Fatal("entry missing after reload")
	}
	if !e.FirstSeen.Equal(fixed) {
		t.Errorf("FirstSeen = %v, want %v", e.FirstSeen, fixed)
	}
}

func TestNew_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	tr, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestAge(t *testing.T) {
	e := Entry{FirstSeen: fixed}
	got := e.Age(fixed.Add(2 * time.Hour))
	if got != 2*time.Hour {
		t.Errorf("Age = %v, want 2h", got)
	}
}

func init() { _ = os.Stderr }
