package config

import (
	"testing"
	"time"
)

func TestApply_OverridesSnapshotPath(t *testing.T) {
	cfg := Default()
	f := &Flags{
		SnapshotPath: "/tmp/override.json",
	}
	result := Apply(cfg, f)
	if result.SnapshotPath != "/tmp/override.json" {
		t.Errorf("expected /tmp/override.json, got %s", result.SnapshotPath)
	}
}

func TestApply_OverridesInterval(t *testing.T) {
	cfg := Default()
	f := &Flags{
		Interval: 15 * time.Second,
	}
	result := Apply(cfg, f)
	if result.ScanInterval != 15*time.Second {
		t.Errorf("expected 15s, got %v", result.ScanInterval)
	}
}

func TestApply_VerboseFlag(t *testing.T) {
	cfg := Default()
	f := &Flags{Verbose: true}
	result := Apply(cfg, f)
	if !result.Verbose {
		t.Error("expected Verbose to be true after Apply")
	}
}

func TestApply_NoOverride(t *testing.T) {
	cfg := Default()
	originalPath := cfg.SnapshotPath
	originalInterval := cfg.ScanInterval
	f := &Flags{} // nothing set
	result := Apply(cfg, f)
	if result.SnapshotPath != originalPath {
		t.Errorf("snapshot path should not change: got %s", result.SnapshotPath)
	}
	if result.ScanInterval != originalInterval {
		t.Errorf("interval should not change: got %v", result.ScanInterval)
	}
}
