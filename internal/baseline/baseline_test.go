package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

var samplePorts = []scanner.Port{
	{Port: 22, Proto: "tcp", Service: "ssh"},
	{Port: 80, Proto: "tcp", Service: "http"},
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	m := baseline.New(path)

	if err := m.Save(samplePorts); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b.Ports) != len(samplePorts) {
		t.Errorf("expected %d ports, got %d", len(samplePorts), len(b.Ports))
	}
	if b.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	m := baseline.New(filepath.Join(dir, "baseline.json"))
	if m.Exists() {
		t.Error("should not exist before Save")
	}
	_ = m.Save(samplePorts)
	if !m.Exists() {
		t.Error("should exist after Save")
	}
}

func TestCompare_Deviation(t *testing.T) {
	b := &baseline.Baseline{
		CapturedAt: time.Now(),
		Ports:      samplePorts,
	}
	current := []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 443, Proto: "tcp", Service: "https"},
	}
	d := baseline.Compare(b, current)
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(d.Added) != 1 || d.Added[0].Port != 443 {
		t.Errorf("unexpected added: %v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed[0].Port != 80 {
		t.Errorf("unexpected removed: %v", d.Removed)
	}
}

func TestCompare_NoDeviation(t *testing.T) {
	b := &baseline.Baseline{CapturedAt: time.Now(), Ports: samplePorts}
	d := baseline.Compare(b, samplePorts)
	if d.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	m := baseline.New(filepath.Join(t.TempDir(), "missing.json"))
	_, err := m.Load()
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
