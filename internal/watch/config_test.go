package watch_test

import (
	"flag"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watch"
)

func TestDefaultConfig(t *testing.T) {
	cfg := watch.DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.SnapshotPath == "" {
		t.Error("expected non-empty default snapshot path")
	}
}

func TestApplyFlags_Interval(t *testing.T) {
	cfg := watch.DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("interval", "30s", "")
	fs.String("snapshot", cfg.SnapshotPath, "")
	_ = fs.Parse([]string{"-interval=5s"})

	watch.ApplyFlags(&cfg, fs)
	if cfg.Interval != 5*time.Second {
		t.Errorf("expected 5s, got %v", cfg.Interval)
	}
}

func TestApplyFlags_Snapshot(t *testing.T) {
	cfg := watch.DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("interval", "30s", "")
	fs.String("snapshot", cfg.SnapshotPath, "")
	_ = fs.Parse([]string{"-snapshot=/tmp/snap.json"})

	watch.ApplyFlags(&cfg, fs)
	if cfg.SnapshotPath != "/tmp/snap.json" {
		t.Errorf("expected /tmp/snap.json, got %s", cfg.Snapshotn
func TestApplyFlags_NoOverride(t *testing.T) {
	cfg := watch.DefaultConfig()
	orig := cfg.Interval
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("interval", "30s", "")
	_ = fs.Parse([]string{})

	watch.ApplyFlags(&cfg, fs)
	if cfg.Interval != orig {
		t.Errorf("interval should not change, got %v", cfg.Interval)
	}
}
