package portmap

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portmap.json")
}

func TestSet_And_Get(t *testing.T) {
	m, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	m.Set(8080, "tcp", "my-api")
	name, ok := m.Get(8080, "tcp")
	if !ok || name != "my-api" {
		t.Fatalf("expected my-api, got %q ok=%v", name, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	m, _ := New(tempPath(t))
	_, ok := m.Get(9999, "tcp")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRemove(t *testing.T) {
	m, _ := New(tempPath(t))
	m.Set(22, "tcp", "ssh-custom")
	m.Remove(22, "tcp")
	_, ok := m.Get(22, "tcp")
	if ok {
		t.Fatal("expected entry removed")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	m, _ := New(path)
	m.Set(443, "tcp", "https-internal")
	m.Set(53, "udp", "dns-local")
	if err := m.Save(); err != nil {
		t.Fatal(err)
	}
	m2, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	if name, ok := m2.Get(443, "tcp"); !ok || name != "https-internal" {
		t.Errorf("443/tcp: got %q ok=%v", name, ok)
	}
	if name, ok := m2.Get(53, "udp"); !ok || name != "dns-local" {
		t.Errorf("53/udp: got %q ok=%v", name, ok)
	}
}

func TestNew_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	m, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("expected non-nil map")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	path := tempPath(t)
	os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := New(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
