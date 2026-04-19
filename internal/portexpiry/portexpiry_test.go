package portexpiry

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "expiry.json")
}

func TestTrack_And_Expired(t *testing.T) {
	tr, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Track(8080, "tcp")

	// not expired yet
	if got := tr.Expired(time.Hour); len(got) != 0 {
		t.Fatalf("expected 0 expired, got %d", len(got))
	}

	// advance time
	tr.now = func() time.Time { return now.Add(2 * time.Hour) }
	if got := tr.Expired(time.Hour); len(got) != 1 {
		t.Fatalf("expected 1 expired, got %d", len(got))
	}
}

func TestEvict_RemovesEntry(t *testing.T) {
	tr, _ := New(tempPath(t))
	tr.Track(443, "tcp")
	tr.Evict(443, "tcp")
	if got := tr.Expired(0); len(got) != 0 {
		t.Fatalf("expected 0 after evict, got %d", len(got))
	}
}

func TestSaveAndLoad(t *testing.T) {
	p := tempPath(t)
	tr, _ := New(p)
	tr.Track(22, "tcp")
	if err := tr.Save(); err != nil {
		t.Fatal(err)
	}
	tr2, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := tr2.entries[key(22, "tcp")]; !ok {
		t.Fatal("entry not persisted")
	}
}

func TestNew_MissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	tr, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.entries) != 0 {
		t.Fatal("expected empty tracker")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	p := tempPath(t)
	os.WriteFile(p, []byte("not json"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
