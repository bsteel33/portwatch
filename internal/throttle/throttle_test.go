package throttle

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCall(t *testing.T) {
	th := New(DefaultConfig())
	if !th.Allow("tcp:80:added") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinCooldown(t *testing.T) {
	base := time.Now()
	th := New(Config{Cooldown: 10 * time.Minute})
	th.now = fixedNow(base)
	th.Allow("tcp:80:added")
	th.now = fixedNow(base.Add(1 * time.Minute))
	if th.Allow("tcp:80:added") {
		t.Fatal("expected call within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldown(t *testing.T) {
	base := time.Now()
	th := New(Config{Cooldown: 5 * time.Minute})
	th.now = fixedNow(base)
	th.Allow("tcp:80:added")
	th.now = fixedNow(base.Add(6 * time.Minute))
	if !th.Allow("tcp:80:added") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestReset(t *testing.T) {
	base := time.Now()
	th := New(Config{Cooldown: 10 * time.Minute})
	th.now = fixedNow(base)
	th.Allow("tcp:443:removed")
	th.Reset("tcp:443:removed")
	if !th.Allow("tcp:443:removed") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll(t *testing.T) {
	base := time.Now()
	th := New(Config{Cooldown: 10 * time.Minute})
	th.now = fixedNow(base)
	th.Allow("a")
	th.Allow("b")
	th.ResetAll()
	if !th.Allow("a") || !th.Allow("b") {
		t.Fatal("expected all keys allowed after ResetAll")
	}
}
