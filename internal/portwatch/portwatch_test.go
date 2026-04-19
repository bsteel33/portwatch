package portwatch

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempSnap(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", cfg.Protocol)
	}
	if cfg.Timeout != 2*time.Second {
		t.Errorf("unexpected timeout %v", cfg.Timeout)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultConfig()
	if msg := Validate(cfg); msg != "" {
		t.Errorf("unexpected error: %s", msg)
	}
}

func TestValidate_BadProto(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Protocol = "icmp"
	if msg := Validate(cfg); msg == "" {
		t.Error("expected validation error for bad protocol")
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{Protocol: "udp", Timeout: 5 * time.Second}
	ApplyFlags(&dst, src)
	if dst.Protocol != "udp" {
		t.Errorf("expected udp, got %s", dst.Protocol)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	origProto := dst.Protocol
	ApplyFlags(&dst, Config{})
	if dst.Protocol != origProto {
		t.Error("protocol should not have changed")
	}
}

func TestRegisterFlags(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)
	if err := fs.Parse([]string{"-proto", "udp"}); err != nil {
		t.Fatal(err)
	}
	if cfg.Protocol != "udp" {
		t.Errorf("expected udp, got %s", cfg.Protocol)
	}
}

func TestNew_MissingSnapshot(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SnapshotPath = filepath.Join(t.TempDir(), "missing.json")
	w, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestClampTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Timeout = 10 * time.Second
	clampTimeout(&cfg, 5*time.Second)
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s, got %v", cfg.Timeout)
	}
}

func init() { os.Unsetenv("CI") }
