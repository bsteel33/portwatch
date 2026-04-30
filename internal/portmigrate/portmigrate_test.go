package portmigrate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeRecord(t *testing.T, dir string, rec Record) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestMigrate_AlreadyCurrent(t *testing.T) {
	m := New()
	payload, _ := json.Marshal([]map[string]interface{}{{"port": 80, "proto": "tcp"}})
	rec := Record{Version: CurrentVersion, Payload: payload}
	out, err := m.Migrate(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Version != CurrentVersion {
		t.Errorf("version = %d, want %d", out.Version, CurrentVersion)
	}
}

func TestMigrate_V1AddsProto(t *testing.T) {
	m := New()
	payload, _ := json.Marshal([]map[string]interface{}{{"port": 443}})
	rec := Record{Version: Version1, Payload: payload}
	out, err := m.Migrate(rec)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if out.Version != CurrentVersion {
		t.Errorf("version = %d, want %d", out.Version, CurrentVersion)
	}
	var ports []map[string]interface{}
	if err := json.Unmarshal(out.Payload, &ports); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if len(ports) == 0 {
		t.Fatal("expected ports in result")
	}
	if ports[0]["proto"] != "tcp" {
		t.Errorf("proto = %v, want tcp", ports[0]["proto"])
	}
}

func TestMigrate_UnknownVersion(t *testing.T) {
	m := New()
	rec := Record{Version: Version(99), Payload: json.RawMessage(`[]`)}
	_, err := m.Migrate(rec)
	if err == nil {
		t.Error("expected error for unknown version")
	}
}

func TestLoadAndMigrate_MissingFile(t *testing.T) {
	m := New()
	_, err := m.LoadAndMigrate("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadAndMigrate_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	payload, _ := json.Marshal([]map[string]interface{}{{"port": 22}})
	rec := Record{Version: Version1, Payload: payload}
	path := writeRecord(t, dir, rec)

	m := New()
	out, err := m.LoadAndMigrate(path)
	if err != nil {
		t.Fatalf("LoadAndMigrate: %v", err)
	}
	if out.Version != CurrentVersion {
		t.Errorf("version = %d, want %d", out.Version, CurrentVersion)
	}
	if out.MigratedAt.IsZero() {
		t.Error("expected MigratedAt to be set")
	}
}
