// Package portevents provides a simple event bus for broadcasting port change
// events to multiple subscribers within portwatch.
package portevents

import (
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// EventType classifies the kind of port event.
type EventType string

const (
	EventAdded   EventType = "added"
	EventRemoved EventType = "removed"
)

// Event carries a single port change notification.
type Event struct {
	Type EventType
	Port snapshot.Port
}

// Handler is a function that receives port events.
type Handler func(Event)

// Bus is a thread-safe event bus that fans out port events to registered handlers.
type Bus struct {
	mu       sync.RWMutex
	handlers []Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{}
}

// Subscribe registers h to receive future events. Returns an unsubscribe func.
func (b *Bus) Subscribe(h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	idx := len(b.handlers)
	b.handlers = append(b.handlers, h)

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		b.handlers[idx] = nil
	}
}

// Publish sends e to all registered (non-nil) handlers.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers))
	copy(handlers, b.handlers)
	b.mu.RUnlock()

	for _, h := range handlers {
		if h != nil {
			h(e)
		}
	}
}

// PublishDiff emits one event per added and removed port from diff.
func (b *Bus) PublishDiff(added, removed []snapshot.Port) {
	for _, p := range added {
		b.Publish(Event{Type: EventAdded, Port: p})
	}
	for _, p := range removed {
		b.Publish(Event{Type: EventRemoved, Port: p})
	}
}
