// Package portmigrate provides utilities for migrating port snapshot data
// between schema versions, applying transformations as needed.
package portmigrate

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Version represents a snapshot schema version.
type Version int

const (
	VersionUnknown Version = 0
	Version1       Version = 1
	Version2       Version = 2
	CurrentVersion         = Version2
)

// Record holds a versioned snapshot payload for migration.
type Record struct {
	Version   Version         `json:"version"`
	MigratedAt time.Time      `json:"migrated_at,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

// Migrator applies sequential version upgrades to snapshot records.
type Migrator struct {
	steps map[Version]MigrateFunc
}

// MigrateFunc transforms a raw JSON payload from one version to the next.
type MigrateFunc func(payload json.RawMessage) (json.RawMessage, error)

// New returns a Migrator with the default registered migration steps.
func New() *Migrator {
	m := &Migrator{steps: make(map[Version]MigrateFunc)}
	m.Register(Version1, migrateV1toV2)
	return m
}

// Register adds a migration step from the given version to the next.
func (m *Migrator) Register(from Version, fn MigrateFunc) {
	m.steps[from] = fn
}

// Migrate upgrades rec to CurrentVersion, applying each step in order.
func (m *Migrator) Migrate(rec Record) (Record, error) {
	for rec.Version < CurrentVersion {
		fn, ok := m.steps[rec.Version]
		if !ok {
			return rec, fmt.Errorf("portmigrate: no migration from version %d", rec.Version)
		}
		next, err := fn(rec.Payload)
		if err != nil {
			return rec, fmt.Errorf("portmigrate: step %d: %w", rec.Version, err)
		}
		rec.Payload = next
		rec.Version++
	}
	rec.MigratedAt = time.Now().UTC()
	return rec, nil
}

// LoadAndMigrate reads a Record from path and migrates it to CurrentVersion.
func (m *Migrator) LoadAndMigrate(path string) (Record, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, fmt.Errorf("portmigrate: read %s: %w", path, err)
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("portmigrate: parse %s: %w", path, err)
	}
	return m.Migrate(rec)
}

// migrateV1toV2 adds a default "proto" field to any port entries missing it.
func migrateV1toV2(payload json.RawMessage) (json.RawMessage, error) {
	var ports []map[string]interface{}
	if err := json.Unmarshal(payload, &ports); err != nil {
		return payload, err
	}
	for _, p := range ports {
		if _, ok := p["proto"]; !ok {
			p["proto"] = "tcp"
		}
	}
	return json.Marshal(ports)
}
