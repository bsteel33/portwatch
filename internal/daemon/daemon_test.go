package daemon

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func tempSnapshot(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-snap-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Close()
	os.Remove(f.Name()) // daemon should create it
	return f.Name()
}

func TestDaemon_RunAndStop(t *testing.T) {
	cfg := config.Default()
	cfg.SnapshotPath = tempSnapshot(t)
	cfg.Interval = 1
	defer os.Remove(cfg.SnapshotPath)

	d := New(cfg)
	stop := make(chan struct{})

	done := make(chan error, 1)
	go func() {
		done <- d.Run(stop)
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not stop in time")
	}
}

func TestDaemon_CreatesSnapshot(t *testing.T) {
	cfg := config.Default()
	cfg.SnapshotPath = tempSnapshot(t)
	cfg.Interval = 60
	defer os.Remove(cfg.SnapshotPath)

	d := New(cfg)
	if err := d.tick(); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	if _, err := os.Stat(cfg.SnapshotPath); os.IsNotExist(err) {
		t.Fatal("snapshot file was not created")
	}
}
