package portcache

import (
	"flag"
	"testing"
	"time"
)

func TestRegisterFlags_SetsDefault(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)

	if err := fs.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	if cfg.TTL != 30*time.Second {
		t.Fatalf("expected 30s default TTL, got %v", cfg.TTL)
	}
}

func TestRegisterFlags_Override(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)

	if err := fs.Parse([]string{"-cache-ttl=10s"}); err != nil {
		t.Fatal(err)
	}
	if cfg.TTL != 10*time.Second {
		t.Fatalf("expected 10s TTL, got %v", cfg.TTL)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	src := Config{TTL: 5 * time.Second}
	ApplyFlags(&dst, src)
	if dst.TTL != 5*time.Second {
		t.Fatalf("expected 5s TTL after apply, got %v", dst.TTL)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	orig := dst.TTL
	ApplyFlags(&dst, Config{})
	if dst.TTL != orig {
		t.Fatalf("expected TTL unchanged, got %v", dst.TTL)
	}
}
