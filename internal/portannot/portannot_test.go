package portannot

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "annot.json")
}

func TestSet_And_Get(t *testing.T) {
	a, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Set(80, "tcp", "web server"); err != nil {
		t.Fatal(err)
	}
	ann, ok := a.Get(80, "tcp")
	if !ok {
		t.Fatal("expected annotation")
	}
	if ann.Note != "web server" {
		t.Errorf("got %q, want %q", ann.Note, "web server")
	}
}

func TestGet_Missing(t *testing.T) {
	a, _ := New(tempPath(t))
	_, ok := a.Get(9999, "tcp")
	if ok {
		t.Error("expected no annotation")
	}
}

func TestRemove(t *testing.T) {
	p := tempPath(t)
	a, _ := New(p)
	_ = a.Set(443, "tcp", "https")
	_ = a.Remove(443, "tcp")
	_, ok := a.Get(443, "tcp")
	if ok {
		t.Error("expected annotation to be removed")
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	a, _ := New(p)
	_ = a.Set(22, "tcp", "ssh")

	a2, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	ann, ok := a2.Get(22, "tcp")
	if !ok {
		t.Fatal("expected annotation after reload")
	}
	if ann.Note != "ssh" {
		t.Errorf("got %q", ann.Note)
	}
}

func TestNew_MissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	a, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.All()) != 0 {
		t.Error("expected empty store")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := New(p)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
