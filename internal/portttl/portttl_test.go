package portttl

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portttl.json")
}

func TestTrack_And_Expired(t *testing.T) {
	tr, err := New(tempPath(t), fixedNow(epoch))
	if err != nil {
		t.Fatal(err)
	}
	tr.Track(8080, "tcp", time.Hour)
	if got := tr.Expired(); len(got) != 0 {
		t.Fatalf("expected no expired entries, got %d", len(got))
	}
}

func TestExpired_AfterTTL(t *testing.T) {
	path := tempPath(t)
	tr, _ := New(path, fixedNow(epoch))
	tr.Track(22, "tcp", time.Minute)

	// advance clock past TTL
	tr.now = fixedNow(epoch.Add(2 * time.Minute))
	exp := tr.Expired()
	if len(exp) != 1 {
		t.Fatalf("expected 1 expired entry, got %d", len(exp))
	}
	if exp[0].Port != 22 {
		t.Errorf("expected port 22, got %d", exp[0].Port)
	}
}

func TestEvict_RemovesEntry(t *testing.T) {
	tr, _ := New(tempPath(t), fixedNow(epoch))
	tr.Track(443, "tcp", time.Hour)
	tr.Evict(443, "tcp")
	if len(tr.entries) != 0 {
		t.Errorf("expected empty entries after evict")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	tr, _ := New(path, fixedNow(epoch))
	tr.Track(80, "tcp", 2*time.Hour)
	if err := tr.Save(); err != nil {
		t.Fatal(err)
	}
	tr2, err := New(path, fixedNow(epoch))
	if err != nil {
		t.Fatal(err)
	}
	if len(tr2.entries) != 1 {
		t.Errorf("expected 1 entry after reload, got %d", len(tr2.entries))
	}
}

func TestNew_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	tr, err := New(path, fixedNow(epoch))
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.entries) != 0 {
		t.Errorf("expected empty tracker for missing file")
	}
	_ = os.Remove(path)
}
