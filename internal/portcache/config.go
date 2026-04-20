package portcache

import (
	"flag"
	"time"
)

// Config controls cache behaviour.
type Config struct {
	// TTL is how long a cached result remains valid.
	TTL time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		TTL: 30 * time.Second,
	}
}

// RegisterFlags registers cache-related CLI flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.DurationVar(&cfg.TTL, "cache-ttl", cfg.TTL,
		"how long to cache port scan results (0 disables caching)")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.TTL != 0 {
		dst.TTL = src.TTL
	}
}
