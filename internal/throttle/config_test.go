package throttle

import (
	"flag"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Cooldown != 5*time.Minute {
		t.Fatalf("expected 5m cooldown, got %v", cfg.Cooldown)
	}
}

func TestRegisterFlags(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)
	if err := fs.Parse([]string{"-throttle=2m"}); err != nil {
		t.Fatal(err)
	}
	if cfg.Cooldown != 2*time.Minute {
		t.Fatalf("expected 2m, got %v", cfg.Cooldown)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	cfg := DefaultConfig()
	ApplyFlags(&cfg, Config{Cooldown: 30 * time.Second})
	if cfg.Cooldown != 30*time.Second {
		t.Fatalf("expected 30s, got %v", cfg.Cooldown)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	cfg := DefaultConfig()
	ApplyFlags(&cfg, Config{})
	if cfg.Cooldown != 5*time.Minute {
		t.Fatalf("expected default 5m, got %v", cfg.Cooldown)
	}
}
