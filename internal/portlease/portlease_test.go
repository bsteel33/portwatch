package portlease

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "leases.json")
}

func fixedNow(ts time.Time) func() time.Time { return func() time.Time { return ts } }

func TestClaim_And_Get(t *testing.T) {
	l, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	l.now = fixedNow(now)

	if err := l.Claim(8080, "tcp", "alice", time.Minute); err != nil {
		t.Fatal(err)
	}
	lease, ok := l.Get(8080, "tcp")
	if !ok {
		t.Fatal("expected lease to exist")
	}
	if lease.Owner != "alice" {
		t.Errorf("expected owner alice, got %s", lease.Owner)
	}
}

func TestGet_Expired(t *testing.T) {
	l, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	base := time.Now()
	l.now = fixedNow(base)
	_ = l.Claim(443, "tcp", "bob", time.Second)

	l.now = fixedNow(base.Add(2 * time.Second))
	_, ok := l.Get(443, "tcp")
	if ok {
		t.Error("expected expired lease to be absent")
	}
}

func TestRelease_RemovesLease(t *testing.T) {
	l, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	_ = l.Claim(22, "tcp", "carol", time.Hour)
	if err := l.Release(22, "tcp"); err != nil {
		t.Fatal(err)
	}
	_, ok := l.Get(22, "tcp")
	if ok {
		t.Error("expected lease to be released")
	}
}

func TestPersistence(t *testing.T) {
	path := tempPath(t)
	l, _ := New(path)
	_ = l.Claim(9090, "udp", "dave", time.Hour)

	l2, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	lease, ok := l2.Get(9090, "udp")
	if !ok {
		t.Fatal("expected persisted lease")
	}
	if lease.Owner != "dave" {
		t.Errorf("expected dave, got %s", lease.Owner)
	}
}

func TestNew_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	_, err := New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestFprint_NoLeases(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, map[string]Lease{}, time.Now())
	if buf.String() != "no active leases\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFprint_WithLeases(t *testing.T) {
	now := time.Now()
	leases := map[string]Lease{
		"80/tcp": {Owner: "eve", ExpiresAt: now.Add(30 * time.Second)},
	}
	var buf bytes.Buffer
	Fprint(&buf, leases, now)
	if !bytes.Contains(buf.Bytes(), []byte("eve")) {
		t.Error("expected owner name in output")
	}
	_ = os.Stdout
}
