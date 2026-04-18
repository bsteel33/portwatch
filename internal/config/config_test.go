package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.ScanInterval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.ScanInterval)
	}
	if cfg.SnapshotPath == "" {
		t.Error("expected non-empty snapshot path")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	cfg := &Config{
		Ports:        []int{22, 80, 443},
		ScanInterval: 30 * time.Second,
		SnapshotPath: "/tmp/snap.json",
		AlertEmail:   "ops@example.com",
		Verbose:      true,
	}

	if err := Save(cfg, tmp.Name()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.AlertEmail != cfg.AlertEmail {
		t.Errorf("AlertEmail mismatch: got %q want %q", loaded.AlertEmail, cfg.AlertEmail)
	}
	if loaded.ScanInterval != cfg.ScanInterval {
		t.Errorf("ScanInterval mismatch: got %v want %v", loaded.ScanInterval, cfg.ScanInterval)
	}
	if len(loaded.Ports) != len(cfg.Ports) {
		t.Errorf("Ports length mismatch: got %d want %d", len(loaded.Ports), len(cfg.Ports))
	}
	if !loaded.Verbose {
		t.Error("expected Verbose to be true")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/portwatch-config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
