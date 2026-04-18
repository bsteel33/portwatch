package snapshot_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func samplePorts() []snapshot.PortInfo {
	return []snapshot.PortInfo{
		{Port: 22, Service: "ssh", Proto: "tcp"},
		{Port: 80, Service: "http", Proto: "tcp"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	f, err := os.CreateTemp("", "snapshot-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	s := snapshot.New(samplePorts())
	if err := s.Save(f.Name()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(f.Name())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Ports) != len(s.Ports) {
		t.Errorf("expected %d ports, got %d", len(s.Ports), len(loaded.Ports))
	}
}

func TestCompare(t *testing.T) {
	old := snapshot.New(samplePorts())
	newPorts := []snapshot.PortInfo{
		{Port: 22, Service: "ssh", Proto: "tcp"},
		{Port: 443, Service: "https", Proto: "tcp"},
	}
	new := snapshot.New(newPorts)

	diff := snapshot.Compare(old, new)

	if !diff.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(diff.Opened) != 1 || diff.Opened[0].Port != 443 {
		t.Errorf("expected port 443 opened, got %v", diff.Opened)
	}
	if len(diff.Closed) != 1 || diff.Closed[0].Port != 80 {
		t.Errorf("expected port 80 closed, got %v", diff.Closed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	old := snapshot.New(samplePorts())
	new := snapshot.New(samplePorts())
	diff := snapshot.Compare(old, new)
	if diff.HasChanges() {
		t.Error("expected no changes")
	}
}
