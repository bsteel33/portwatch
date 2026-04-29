package portrelay_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portrelay"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return snapshot.New([]scanner.Port{
		{Port: 80, Proto: "tcp", Service: "http"},
	})
}

// errDestination always returns an error.
type errDestination struct{ name string }

func (e *errDestination) Name() string { return e.name }
func (e *errDestination) Send(_ *snapshot.Snapshot) error {
	return errors.New("send failed")
}

func TestForward_WriterDestination(t *testing.T) {
	var buf bytes.Buffer
	errBuf := &bytes.Buffer{}
	relay := portrelay.New(errBuf)
	relay.Register(portrelay.NewWriterDestination("buf", &buf))

	snap := makeSnap()
	relay.Forward(snap)

	if buf.Len() == 0 {
		t.Fatal("expected output in writer destination")
	}
	var got snapshot.Snapshot
	if err := json.NewDecoder(&buf).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestForward_ErrorDestination_LogsAndContinues(t *testing.T) {
	var errBuf bytes.Buffer
	relay := portrelay.New(&errBuf)
	relay.Register(&errDestination{name: "bad"})

	var goodBuf bytes.Buffer
	relay.Register(portrelay.NewWriterDestination("good", &goodBuf))

	relay.Forward(makeSnap())

	if errBuf.Len() == 0 {
		t.Error("expected error logged for failing destination")
	}
	if goodBuf.Len() == 0 {
		t.Error("expected good destination to still receive snapshot")
	}
}

func TestLen_ReturnsRegisteredCount(t *testing.T) {
	relay := portrelay.New(&bytes.Buffer{})
	if relay.Len() != 0 {
		t.Fatalf("expected 0, got %d", relay.Len())
	}
	relay.Register(portrelay.NewWriterDestination("a", &bytes.Buffer{}))
	relay.Register(portrelay.NewWriterDestination("b", &bytes.Buffer{}))
	if relay.Len() != 2 {
		t.Fatalf("expected 2, got %d", relay.Len())
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := portrelay.DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.BufferSize <= 0 {
		t.Errorf("expected positive BufferSize, got %d", cfg.BufferSize)
	}
}

func TestHTTPDestination_Name(t *testing.T) {
	d := portrelay.NewHTTPDestination("remote", "http://localhost:9999", time.Second)
	if d.Name() != "remote" {
		t.Errorf("expected 'remote', got %q", d.Name())
	}
}
