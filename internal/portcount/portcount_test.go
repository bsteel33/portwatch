package portcount

import "testing"

func TestUpdate_BelowThreshold(t *testing.T) {
	c := New(Config{MaxPorts: 10})
	msg := c.Update(5)
	if msg != "" {
		t.Fatalf("expected no alert, got %q", msg)
	}
	if c.Current() != 5 {
		t.Fatalf("expected current=5, got %d", c.Current())
	}
}

func TestUpdate_ExceedsThreshold(t *testing.T) {
	c := New(Config{MaxPorts: 10})
	msg := c.Update(15)
	if msg == "" {
		t.Fatal("expected alert message, got empty string")
	}
}

func TestUpdate_NoThreshold(t *testing.T) {
	c := New(Config{MaxPorts: 0})
	msg := c.Update(1000)
	if msg != "" {
		t.Fatalf("expected no alert when MaxPorts=0, got %q", msg)
	}
}

func TestPeak_TracksHighest(t *testing.T) {
	c := New(DefaultConfig())
	c.Update(3)
	c.Update(7)
	c.Update(4)
	if c.Peak() != 7 {
		t.Fatalf("expected peak=7, got %d", c.Peak())
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	c := New(DefaultConfig())
	c.Update(8)
	c.Reset()
	if c.Current() != 0 {
		t.Fatalf("expected current=0 after reset, got %d", c.Current())
	}
	if c.Peak() != 0 {
		t.Fatalf("expected peak=0 after reset, got %d", c.Peak())
	}
}

func TestUpdate_AtThreshold_NoAlert(t *testing.T) {
	c := New(Config{MaxPorts: 10})
	msg := c.Update(10)
	if msg != "" {
		t.Fatalf("expected no alert at exact threshold, got %q", msg)
	}
}
