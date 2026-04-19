package portlabel

import (
	"testing"
)

func seedLabels() []Label {
	return []Label{
		{Port: 80, Proto: "tcp", Name: "http"},
		{Port: 443, Proto: "tcp", Name: "https"},
	}
}

func TestGet_Found(t *testing.T) {
	l := New(seedLabels())
	name, ok := l.Get(80, "tcp")
	if !ok || name != "http" {
		t.Fatalf("expected http, got %q ok=%v", name, ok)
	}
}

func TestGet_NotFound(t *testing.T) {
	l := New(seedLabels())
	_, ok := l.Get(22, "tcp")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSet_Overwrite(t *testing.T) {
	l := New(seedLabels())
	l.Set(80, "tcp", "web")
	name, _ := l.Get(80, "tcp")
	if name != "web" {
		t.Fatalf("expected web, got %q", name)
	}
}

func TestRemove(t *testing.T) {
	l := New(seedLabels())
	l.Remove(80, "tcp")
	_, ok := l.Get(80, "tcp")
	if ok {
		t.Fatal("expected label to be removed")
	}
}

func TestAll_Count(t *testing.T) {
	l := New(seedLabels())
	if got := len(l.All()); got != 2 {
		t.Fatalf("expected 2 labels, got %d", got)
	}
}

func TestParseRules_Valid(t *testing.T) {
	cfg := Config{Rules: []string{"8080/tcp=proxy", "53/udp=dns"}}
	labels, err := ParseRules(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2, got %d", len(labels))
	}
}

func TestParseRules_Invalid(t *testing.T) {
	cfg := Config{Rules: []string{"badrule"}}
	_, err := ParseRules(cfg)
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}
