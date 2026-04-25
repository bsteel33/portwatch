package portevents

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makePort(num int, proto string) snapshot.Port {
	return snapshot.Port{Port: num, Proto: proto}
}

func TestPublish_DeliverToSubscriber(t *testing.T) {
	bus := New()
	var got []Event
	bus.Subscribe(func(e Event) { got = append(got, e) })

	bus.Publish(Event{Type: EventAdded, Port: makePort(80, "tcp")})

	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if got[0].Type != EventAdded {
		t.Errorf("expected EventAdded, got %q", got[0].Type)
	}
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	bus := New()
	var mu sync.Mutex
	count := 0
	inc := func(_ Event) { mu.Lock(); count++; mu.Unlock() }

	bus.Subscribe(inc)
	bus.Subscribe(inc)
	bus.Publish(Event{Type: EventRemoved, Port: makePort(22, "tcp")})

	if count != 2 {
		t.Errorf("expected 2 handler calls, got %d", count)
	}
}

func TestUnsubscribe_StopsDelivery(t *testing.T) {
	bus := New()
	var got []Event
	unsub := bus.Subscribe(func(e Event) { got = append(got, e) })

	unsub()
	bus.Publish(Event{Type: EventAdded, Port: makePort(443, "tcp")})

	if len(got) != 0 {
		t.Errorf("expected no events after unsubscribe, got %d", len(got))
	}
}

func TestPublishDiff_EmitsCorrectTypes(t *testing.T) {
	bus := New()
	var added, removed int
	bus.Subscribe(func(e Event) {
		switch e.Type {
		case EventAdded:
			added++
		case EventRemoved:
			removed++
		}
	})

	bus.PublishDiff(
		[]snapshot.Port{makePort(80, "tcp"), makePort(443, "tcp")},
		[]snapshot.Port{makePort(22, "tcp")},
	)

	if added != 2 {
		t.Errorf("expected 2 added events, got %d", added)
	}
	if removed != 1 {
		t.Errorf("expected 1 removed event, got %d", removed)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.BufferSize != 0 {
		t.Errorf("expected BufferSize 0, got %d", cfg.BufferSize)
	}
}

func TestApplyFlags_Override(t *testing.T) {
	dst := DefaultConfig()
	ApplyFlags(&dst, Config{BufferSize: 32})
	if dst.BufferSize != 32 {
		t.Errorf("expected 32, got %d", dst.BufferSize)
	}
}

func TestApplyFlags_NoOverride(t *testing.T) {
	dst := DefaultConfig()
	ApplyFlags(&dst, Config{BufferSize: 0})
	if dst.BufferSize != 0 {
		t.Errorf("expected 0 (no change), got %d", dst.BufferSize)
	}
}
