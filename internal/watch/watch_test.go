package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/watch"
)

func TestWatcher_CallsOnChange(t *testing.T) {
	s := scanner.New([]int{}, time.Millisecond*100)
	w := watch.New(s, time.Millisecond*50)

	changed := make(chan snapshot.Diff, 1)
	w.OnChange = func(d snapshot.Diff) {
		changed <- d
	}

	dir := t.TempDir()
	snapshotPath := dir + "/snap.json"

	// seed a snapshot with a port that won't be open
	initial := snapshot.New([]snapshot.Port{
		{Port: 19999, Proto: "tcp", Service: "unknown"},
	})
	if err := initial.Save(snapshotPath); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go w.Run(ctx, snapshotPath) //nolint:errcheck

	select {
	case diff := <-changed:
		if len(diff.Removed) == 0 {
			t.Errorf("expected removed ports in diff, got none")
		}
	case <-time.After(800 * time.Millisecond):
		t.Error("timed out waiting for OnChange to be called")
	}
}

func TestWatcher_NoChangeNoCallback(t *testing.T) {
	s := scanner.New([]int{}, time.Millisecond*100)
	w := watch.New(s, time.Millisecond*50)

	called := 0
	w.OnChange = func(snapshot.Diff) { called++ }

	dir := t.TempDir()
	snapshotPath := dir + "/snap.json"

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	go w.Run(ctx, snapshotPath) //nolint:errcheck

	<-ctx.Done()
	// With no previous snapshot and no open ports, diff should be empty after first scan
	if called > 1 {
		t.Errorf("OnChange called %d times, expected at most 1", called)
	}
}
