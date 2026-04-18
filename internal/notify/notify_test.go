package notify

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(added, removed []snapshot.PortEntry) snapshot.Diff {
	return snapshot.Diff{Added: added, Removed: removed}
}

func TestNotify_Stdout_Added(t *testing.T) {
	n := New(Config{Channel: ChannelStdout})
	var buf bytes.Buffer
	n.out = &buf

	diff := makeDiff([]snapshot.PortEntry{{Port: 8080, Proto: "tcp", Service: "http-alt"}}, nil)
	if err := n.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ 8080/tcp") {
		t.Errorf("expected added port in output, got: %s", out)
	}
}

func TestNotify_Stdout_Removed(t *testing.T) {
	n := New(Config{Channel: ChannelStdout})
	var buf bytes.Buffer
	n.out = &buf

	diff := makeDiff(nil, []snapshot.PortEntry{{Port: 22, Proto: "tcp", Service: "ssh"}})
	if err := n.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "- 22/tcp") {
		t.Errorf("expected removed port in output, got: %s", out)
	}
}

func TestNotify_NoChanges(t *testing.T) {
	n := New(Config{Channel: ChannelStdout})
	var buf bytes.Buffer
	n.out = &buf

	if err := n.Notify(makeDiff(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff")
	}
}

func TestNotify_Exec_EmptyCmd(t *testing.T) {
	n := New(Config{Channel: ChannelExec, ExecCmd: ""})
	diff := makeDiff([]snapshot.PortEntry{{Port: 9090, Proto: "tcp", Service: "unknown"}}, nil)
	if err := n.Notify(diff); err == nil {
		t.Error("expected error for empty exec command")
	}
}

func TestBuildMessage(t *testing.T) {
	diff := makeDiff(
		[]snapshot.PortEntry{{Port: 443, Proto: "tcp", Service: "https"}},
		[]snapshot.PortEntry{{Port: 80, Proto: "tcp", Service: "http"}},
	)
	msg := buildMessage(diff)
	if !strings.Contains(msg, "+ 443/tcp") {
		t.Errorf("missing added entry in message")
	}
	if !strings.Contains(msg, "- 80/tcp") {
		t.Errorf("missing removed entry in message")
	}
}
