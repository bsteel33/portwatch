package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(added, removed []snapshot.PortEntry) snapshot.Diff {
	return snapshot.Diff{Added: added, Removed: removed}
}

func TestNotify_Added(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diff := makeDiff(
		[]snapshot.PortEntry{{Port: 8080, Proto: "tcp", Service: "http-alt"}},
		nil,
	)

	alerts := n.Notify(diff)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != LevelWarn {
		t.Errorf("expected WARN level, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "port opened") {
		t.Errorf("output missing 'port opened': %s", buf.String())
	}
}

func TestNotify_Removed(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diff := makeDiff(
		nil,
		[]snapshot.PortEntry{{Port: 22, Proto: "tcp", Service: "ssh"}},
	)

	alerts := n.Notify(diff)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != LevelInfo {
		t.Errorf("expected INFO level, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "port closed") {
		t.Errorf("output missing 'port closed': %s", buf.String())
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	alerts := n.Notify(makeDiff(nil, nil))
	if len(alerts) != 0 {
		t.Errorf("expected no alerts, got %d", len(alerts))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}
