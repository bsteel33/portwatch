package suppress

import (
	"flag"
	"time"
)

// Config holds configuration for the Suppressor.
type Config struct {
	TTL time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		TTL: 10 * time.Minute,
	}
}

// RegisterFlags registers suppress-related flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.TTL, "suppress-ttl", cfg.TTL, "duration to suppress repeated alerts for the same port")
}

// ApplyFlags overwrites cfg fields from src when they differ from the default.
func ApplyFlags(cfg *Config, src Config) {
	def := DefaultConfig()
	if src.TTL != def.TTL {
		cfg.TTL = src.TTL
	}
}
