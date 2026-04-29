package portage

import (
	"flag"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Path == "" {
		t.Fatal("expected non-empty default Path")
	}
	if cfg.MaxAge <= 0 {
		t.Fatal("expected positive default MaxAge")
	}
}

func TestRegisterFlags(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)

	if err := fs.Parse([]string{"-portage.path", "/tmp/pa.json", "-portage.max-age", "48h"}); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.Path != "/tmp/pa.json" {
		t.Errorf("Path = %q, want /tmp/pa.json", cfg.Path)
	}
	if cfg.MaxAge != 48*time.Hour {
		t.Errorf("MaxAge = %v, want 48h", cfg.MaxAge)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{Path: "/var/portage.json", MaxAge: 24 * time.Hour}
	ApplyFlags(&dst, src)
	if dst.Path != src.Path {
		t.Errorf("Path = %q, want %q", dst.Path, src.Path)
	}
	if dst.MaxAge != src.MaxAge {
		t.Errorf("MaxAge = %v, want %v", dst.MaxAge, src.MaxAge)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	origPath := dst.Path
	origMaxAge := dst.MaxAge
	ApplyFlags(&dst, Config{})
	if dst.Path != origPath {
		t.Errorf("Path changed unexpectedly: %q", dst.Path)
	}
	if dst.MaxAge != origMaxAge {
		t.Errorf("MaxAge changed unexpectedly: %v", dst.MaxAge)
	}
}
