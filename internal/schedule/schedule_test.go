package schedule

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestRun_CallsJobImmediately(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 10 * time.Millisecond
	cfg.DelayFirst = false
	sched := New(cfg)

	var calls int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sched.Run(ctx, func() { //nolint
		atomic.AddInt32(&calls, 1)
	})

	time.Sleep(5 * time.Millisecond)
	if atomic.LoadInt32(&calls) < 1 {
		t.Fatal("expected immediate first call")
	}
}

func TestRun_TicksRepeatedly(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 20 * time.Millisecond
	cfg.DelayFirst = true
	sched := New(cfg)

	var calls int32
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	sched.Run(ctx, func() { //nolint
		atomic.AddInt32(&calls, 1)
	})

	got := atomic.LoadInt32(&calls)
	if got < 2 || got > 5 {
		t.Fatalf("expected ~3 ticks, got %d", got)
	}
}

func TestRun_StopsOnCancel(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 5 * time.Millisecond
	sched := New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sched.Run(ctx, func() {})
	if err == nil {
		t.Fatal("expected context error")
	}
}
