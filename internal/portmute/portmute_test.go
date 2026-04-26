package portmute

import (
	"bytes"
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func newTestMuter() *Muter {
	m := New()
	m.now = func() time.Time { return fixedNow }
	return m
}

func TestMute_And_IsMuted(t *testing.T) {
	m := newTestMuter()
	m.Mute(443, "tcp", 10*time.Minute, "test")
	if !m.IsMuted(443, "tcp") {
		t.Fatal("expected port to be muted")
	}
}

func TestIsMuted_Expired(t *testing.T) {
	m := newTestMuter()
	m.Mute(443, "tcp", 5*time.Minute, "")
	// advance clock past expiry
	m.now = func() time.Time { return fixedNow.Add(10 * time.Minute) }
	if m.IsMuted(443, "tcp") {
		t.Fatal("expected mute to have expired")
	}
}

func TestUnmute_ClearsEntry(t *testing.T) {
	m := newTestMuter()
	m.Mute(80, "tcp", time.Hour, "")
	m.Unmute(80, "tcp")
	if m.IsMuted(80, "tcp") {
		t.Fatal("expected port to be unmuted")
	}
}

func TestGet_Found(t *testing.T) {
	m := newTestMuter()
	m.Mute(22, "tcp", time.Hour, "ssh maintenance")
	e, ok := m.Get(22, "tcp")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Reason != "ssh maintenance" {
		t.Fatalf("unexpected reason: %s", e.Reason)
	}
}

func TestGet_Missing(t *testing.T) {
	m := newTestMuter()
	_, ok := m.Get(9999, "udp")
	if ok {
		t.Fatal("expected no entry for unknown port")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	m := newTestMuter()
	m.Mute(80, "tcp", time.Hour, "")
	m.Mute(443, "tcp", time.Hour, "")
	m.Reset()
	if m.IsMuted(80, "tcp") || m.IsMuted(443, "tcp") {
		t.Fatal("expected all mutes cleared after reset")
	}
}

func TestFprint_NoMutes(t *testing.T) {
	m := newTestMuter()
	var buf bytes.Buffer
	Fprint(&buf, m)
	if buf.String() != "no muted ports\n" {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestFprint_ShowsActiveEntries(t *testing.T) {
	m := newTestMuter()
	m.Mute(8080, "tcp", time.Hour, "deploy")
	var buf bytes.Buffer
	Fprint(&buf, m)
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}
	if !bytes.Contains([]byte(out), []byte("tcp:8080")) {
		t.Fatalf("expected port key in output, got: %s", out)
	}
}
