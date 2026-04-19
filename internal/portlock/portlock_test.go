package portlock

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "locks.json")
}

func TestAdd_And_Unlocked(t *testing.T) {
	l, err := New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Add(22, "tcp", "ssh"); err != nil {
		t.Fatal(err)
	}
	ports := []scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
	}
	unlocked := l.Unlocked(ports)
	if len(unlocked) != 1 || unlocked[0].Port != 80 {
		t.Fatalf("expected only port 80 unlocked, got %v", unlocked)
	}
}

func TestRemove(t *testing.T) {
	p := tempPath(t)
	l, _ := New(p)
	_ = l.Add(443, "tcp", "https")
	_ = l.Remove(443, "tcp")
	ports := []scanner.Port{{Port: 443, Proto: "tcp"}}
	unlocked := l.Unlocked(ports)
	if len(unlocked) != 1 {
		t.Fatal("expected port 443 to be unlocked after removal")
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	l, _ := New(p)
	_ = l.Add(8080, "tcp", "dev")

	l2, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	ports := []scanner.Port{{Port: 8080, Proto: "tcp"}, {Port: 9090, Proto: "tcp"}}
	unlocked := l2.Unlocked(ports)
	if len(unlocked) != 1 || unlocked[0].Port != 9090 {
		t.Fatalf("expected only 9090 unlocked after reload, got %v", unlocked)
	}
}

func TestNew_MissingFile(t *testing.T) {
	l, err := New(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatal(err)
	}
	if l == nil {
		t.Fatal("expected non-nil locker")
	}
}

func TestUnlocked_AllLocked(t *testing.T) {
	l, _ := New(tempPath(t))
	_ = l.Add(22, "tcp", "")
	_ = l.Add(80, "tcp", "")
	ports := []scanner.Port{{Port: 22, Proto: "tcp"}, {Port: 80, Proto: "tcp"}}
	if got := l.Unlocked(ports); len(got) != 0 {
		t.Fatalf("expected no unlocked ports, got %v", got)
	}
}

func init() {
	_ = os.Stderr
}
