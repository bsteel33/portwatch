package throttle

import (
	"flag"
	"time"
)

// Config holds throttle configuration.
type Config struct {
	Cooldown time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Cooldown: 5 * time.Minute,
	}
}

// RegisterFlags registers throttle flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.Cooldown, "throttle", cfg.Cooldown, "minimum duration between repeated alerts for the same port")
}

// ApplyFlags overrides cfg fields from src when they differ from the zero value.
func ApplyFlags(cfg *Config, src Config) {
	if src.Cooldown != 0 {
		cfg.Cooldown = src.Cooldown
	}
}
