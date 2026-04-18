package ratelimit

import (
	"flag"
	"time"
)

// Config holds rate limiter settings.
type Config struct {
	MaxEvents int
	Window    time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxEvents: 10,
		Window:    time.Minute,
	}
}

// RegisterFlags registers rate limit flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.MaxEvents, "ratelimit-max", cfg.MaxEvents, "maximum events allowed per window")
	fs.DurationVar(&cfg.Window, "ratelimit-window", cfg.Window, "time window for rate limiting")
}

// ApplyFlags overrides cfg fields from src if they differ from defaults.
func ApplyFlags(cfg *Config, src Config) {
	def := DefaultConfig()
	if src.MaxEvents != def.MaxEvents {
		cfg.MaxEvents = src.MaxEvents
	}
	if src.Window != def.Window {
		cfg.Window = src.Window
	}
}
