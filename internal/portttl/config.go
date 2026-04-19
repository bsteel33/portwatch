package portttl

import (
	"flag"
	"time"
)

// Config holds portttl configuration.
type Config struct {
	Path       string
	DefaultTTL time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:       "portttl.json",
		DefaultTTL: 24 * time.Hour,
	}
}

// RegisterFlags registers portttl flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "ttl-path", cfg.Path, "path to portttl state file")
	fs.DurationVar(&cfg.DefaultTTL, "ttl-default", cfg.DefaultTTL, "default TTL for tracked ports")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.DefaultTTL != 0 {
		dst.DefaultTTL = src.DefaultTTL
	}
}
