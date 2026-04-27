// Package portbatch provides batched processing of port sets with configurable
// batch size and flush behaviour.
package portbatch

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Batch holds a slice of ports collected during one flush window.
type Batch struct {
	Ports     []scanner.Port
	Collected time.Time
}

// Processor is a function that receives a flushed batch.
type Processor func(Batch)

// Batcher accumulates ports and flushes them in batches.
type Batcher struct {
	mu        sync.Mutex
	cfg       Config
	buf       []scanner.Port
	processor Processor
}

// New creates a Batcher with the given config and processor callback.
func New(cfg Config, p Processor) *Batcher {
	return &Batcher{
		cfg:       cfg,
		processor: p,
	}
}

// Add appends a port to the internal buffer. If the buffer reaches BatchSize
// it is flushed immediately.
func (b *Batcher) Add(p scanner.Port) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf = append(b.buf, p)
	if b.cfg.BatchSize > 0 && len(b.buf) >= b.cfg.BatchSize {
		b.flush()
	}
}

// AddAll appends multiple ports and flushes when the batch size is reached.
func (b *Batcher) AddAll(ports []scanner.Port) {
	for _, p := range ports {
		b.Add(p)
	}
}

// Flush forces an immediate flush of the current buffer, even if it is not full.
func (b *Batcher) Flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flush()
}

// Len returns the number of ports currently buffered.
func (b *Batcher) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.buf)
}

// flush drains the buffer and calls the processor. Must be called with mu held.
func (b *Batcher) flush() {
	if len(b.buf) == 0 {
		return
	}
	batch := Batch{
		Ports:     make([]scanner.Port, len(b.buf)),
		Collected: time.Now(),
	}
	copy(batch.Ports, b.buf)
	b.buf = b.buf[:0]
	if b.processor != nil {
		b.processor(batch)
	}
}
