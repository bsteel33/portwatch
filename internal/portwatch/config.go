package portwatch

import (
	"errors"
	"flag"
	"strings"
	"time"
)

// Rule describes a single watch condition.
type Rule struct {
	Name    string
	Port    int
	Proto   string
	Timeout time.Duration
}

// Config holds watcher configuration.
type Config struct {
	SnapshotPath string
	Timeout      time.Duration
	Proto        string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		SnapshotPath: "portwatch.json",
		Timeout:      5 * time.Second,
		Proto:        "tcp",
	}
}

// RegisterFlags registers config flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.SnapshotPath, "portwatch.snapshot", cfg.SnapshotPath, "snapshot file path")
	fs.DurationVar(&cfg.Timeout, "portwatch.timeout", cfg.Timeout, "per-port connection timeout")
	fs.StringVar(&cfg.Proto, "portwatch.proto", cfg.Proto, "default protocol filter (tcp/udp)")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.SnapshotPath != "" {
		dst.SnapshotPath = src.SnapshotPath
	}
	if src.Timeout > 0 {
		dst.Timeout = src.Timeout
	}
	if src.Proto != "" {
		dst.Proto = src.Proto
	}
}

// Validate returns an error if cfg is invalid.
func Validate(cfg Config) error {
	if cfg.SnapshotPath == "" {
		return errors.New("portwatch: snapshot path must not be empty")
	}
	p := strings.ToLower(cfg.Proto)
	if p != "tcp" && p != "udp" && p != "" {
		return errors.New("portwatch: proto must be tcp or udp")
	}
	return clampTimeout(cfg)
}

func clampTimeout(cfg Config) error {
	if cfg.Timeout < 0 {
		return errors.New("portwatch: timeout must be non-negative")
	}
	return nil
}
