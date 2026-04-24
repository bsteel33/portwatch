package portretry

import (
	"errors"
	"testing"
	"time"
)

func TestRun_SucceedsFirstAttempt(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Delay = 0
	r := New(cfg)
	calls := 0
	err := r.Run(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRun_RetriesOnFailure(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 3
	cfg.Delay = 0
	r := New(cfg)
	calls := 0
	err := r.Run(func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRun_AllAttemptsFail(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 2
	cfg.Delay = 0
	r := New(cfg)
	err := r.Run(func() error {
		return errors.New("permanent failure")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_SleepsBetweenAttempts(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxAttempts = 3
	cfg.Delay = 100 * time.Millisecond
	r := New(cfg)
	slept := 0
	r.sleep = func(d time.Duration) { slept++ }
	_ = r.Run(func() error { return errors.New("fail") })
	if slept != 2 {
		t.Fatalf("expected 2 sleeps, got %d", slept)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	cfg := DefaultConfig()
	ApplyFlags(&cfg, 5, 2*time.Second)
	if cfg.MaxAttempts != 5 {
		t.Fatalf("expected 5, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != 2*time.Second {
		t.Fatalf("expected 2s, got %v", cfg.Delay)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	cfg := DefaultConfig()
	orig := cfg
	ApplyFlags(&cfg, 0, 0)
	if cfg != orig {
		t.Fatal("config should not change when flags are zero")
	}
}
