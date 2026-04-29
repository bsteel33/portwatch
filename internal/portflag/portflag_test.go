package portflag

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portflag.json")
}

func TestSet_And_Has(t *testing.T) {
	f, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Set(80, "tcp", "watched"); err != nil {
		t.Fatal(err)
	}
	if !f.Has(80, "tcp", "watched") {
		t.Error("expected flag 'watched' to be set")
	}
	if f.Has(80, "tcp", "ignored") {
		t.Error("unexpected flag 'ignored'")
	}
}

func TestUnset_RemovesFlag(t *testing.T) {
	f, _ := New(tempPath(t))
	_ = f.Set(443, "tcp", "reviewed")
	_ = f.Unset(443, "tcp", "reviewed")
	if f.Has(443, "tcp", "reviewed") {
		t.Error("expected flag to be removed")
	}
}

func TestFlags_ReturnsAll(t *testing.T) {
	f, _ := New(tempPath(t))
	_ = f.Set(22, "tcp", "watched")
	_ = f.Set(22, "tcp", "critical")

	got := f.Flags(22, "tcp")
	sort.Strings(got)
	if len(got) != 2 || got[0] != "critical" || got[1] != "watched" {
		t.Errorf("unexpected flags: %v", got)
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	f1, _ := New(p)
	_ = f1.Set(8080, "tcp", "ignored")

	f2, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	if !f2.Has(8080, "tcp", "ignored") {
		t.Error("expected persisted flag to be present after reload")
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := New(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestNew_CorruptFile(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json{"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
