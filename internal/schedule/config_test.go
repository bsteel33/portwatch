package schedule

import (
	"flag"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Fatalf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.DelayFirst {
		t.Fatal("expected DelayFirst false")
	}
}

func TestRegisterFlags(t *testing.T) {
	cfg := DefaultConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs, &cfg)
	fs.Parse([]string{"--schedule.interval=1m", "--schedule.delay-first"})
	if cfg.Interval != time.Minute {
		t.Fatalf("expected 1m, got %v", cfg.Interval)
	}
	if !cfg.DelayFirst {
		t.Fatal("expected DelayFirst true")
	}
}

func TestApplyFlags_Override(t *testing.T) {
	cfg := DefaultConfig()
	ApplyFlags(&cfg, Config{Interval: 5 * time.Minute, DelayFirst: true})
	if cfg.Interval != 5*time.Minute {
		t.Fatalf("unexpected interval %v", cfg.Interval)
	}
	if !cfg.DelayFirst {
		t.Fatal("expected DelayFirst true")
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	cfg := DefaultConfig()
	ApplyFlags(&cfg, Config{})
	if cfg.Interval != 30*time.Second {
		t.Fatalf("interval should not change, got %v", cfg.Interval)
	}
}
