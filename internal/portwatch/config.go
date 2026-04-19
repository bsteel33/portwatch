package portwatch

import (
	"flag"
	"time"
)

// RegisterFlags registers Watcher flags onto fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.SnapshotPath, "snapshot", cfg.SnapshotPath, "path to snapshot file")
	fs.StringVar(&cfg.Protocol, "proto", cfg.Protocol, "protocol to scan (tcp/udp)")
	fs.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "per-port dial timeout")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.SnapshotPath != "" {
		dst.SnapshotPath = src.SnapshotPath
	}
	if src.Protocol != "" {
		dst.Protocol = src.Protocol
	}
	if src.Timeout > 0 {
		dst.Timeout = src.Timeout
	}
}

// Validate returns an error string if cfg is invalid, or empty string if ok.
func Validate(cfg Config) string {
	if cfg.SnapshotPath == "" {
		return "snapshot path must not be empty"
	}
	if cfg.Protocol != "tcp" && cfg.Protocol != "udp" {
		return "protocol must be tcp or udp"
	}
	if cfg.Timeout <= 0 {
		return "timeout must be positive"
	}
	return ""
}

// clampTimeout ensures timeout does not exceed max.
func clampTimeout(cfg *Config, max time.Duration) {
	if cfg.Timeout > max {
		cfg.Timeout = max
	}
}
