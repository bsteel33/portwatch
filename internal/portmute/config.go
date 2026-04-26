package portmute

import (
	"flag"
	"time"
)

// Config holds configuration for the portmute module.
type Config struct {
	// DefaultDuration is used when no explicit duration is specified via CLI.
	DefaultDuration time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultDuration: 30 * time.Minute,
	}
}

// RegisterFlags registers portmute flags on the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.DefaultDuration, "mute-duration", cfg.DefaultDuration,
		"default duration to mute a port when no explicit duration is given")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.DefaultDuration != 0 {
		dst.DefaultDuration = src.DefaultDuration
	}
}
