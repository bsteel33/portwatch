package portbatch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Port: 8000 + i, Proto: "tcp", Service: "test"}
	}
	return ports
}

func TestAdd_FlushesOnBatchSize(t *testing.T) {
	cfg := Config{BatchSize: 3}
	var got []Batch
	b := New(cfg, func(batch Batch) { got = append(got, batch) })

	for _, p := range makePorts(6) {
		b.Add(p)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 flushes, got %d", len(got))
	}
	for _, batch := range got {
		if len(batch.Ports) != 3 {
			t.Errorf("expected batch of 3, got %d", len(batch.Ports))
		}
	}
}

func TestFlush_ExplicitFlush(t *testing.T) {
	cfg := Config{BatchSize: 100}
	var got []Batch
	b := New(cfg, func(batch Batch) { got = append(got, batch) })

	b.AddAll(makePorts(5))
	if len(got) != 0 {
		t.Fatal("expected no automatic flush yet")
	}

	b.Flush()
	if len(got) != 1 {
		t.Fatalf("expected 1 flush after explicit Flush(), got %d", len(got))
	}
	if len(got[0].Ports) != 5 {
		t.Errorf("expected 5 ports in batch, got %d", len(got[0].Ports))
	}
}

func TestFlush_EmptyBufferIsNoop(t *testing.T) {
	cfg := DefaultConfig()
	var got []Batch
	b := New(cfg, func(batch Batch) { got = append(got, batch) })

	b.Flush()
	if len(got) != 0 {
		t.Fatalf("expected no flush for empty buffer, got %d", len(got))
	}
}

func TestLen_ReturnsBufferedCount(t *testing.T) {
	cfg := Config{BatchSize: 100}
	b := New(cfg, nil)

	b.AddAll(makePorts(7))
	if b.Len() != 7 {
		t.Errorf("expected Len()=7, got %d", b.Len())
	}
}

func TestBatch_CollectedTimestamp(t *testing.T) {
	cfg := Config{BatchSize: 1}
	var got []Batch
	b := New(cfg, func(batch Batch) { got = append(got, batch) })

	before := time.Now()
	b.Add(makePorts(1)[0])
	after := time.Now()

	if len(got) != 1 {
		t.Fatal("expected one batch")
	}
	if got[0].Collected.Before(before) || got[0].Collected.After(after) {
		t.Errorf("batch timestamp %v out of expected range [%v, %v]",
			got[0].Collected, before, after)
	}
}
