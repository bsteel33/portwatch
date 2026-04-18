package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Ports to scan; empty means use default range.
	Ports []int `json:"ports"`
	// ScanInterval is how often to re-scan.
	ScanInterval time.Duration `json:"scan_interval"`
	// SnapshotPath is where the snapshot file is stored.
	SnapshotPath string `json:"snapshot_path"`
	// AlertEmail is an optional address to send alerts to.
	AlertEmail string `json:"alert_email,omitempty"`
	// Verbose enables extra logging.
	Verbose bool `json:"verbose"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Ports:        []int{},
		ScanInterval: 60 * time.Second,
		SnapshotPath: "/var/lib/portwatch/snapshot.json",
		Verbose:      false,
	}
}

// Load reads a JSON config file from path and merges it over defaults.
func Load(path string) (*Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config as JSON to path.
func Save(cfg *Config, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
