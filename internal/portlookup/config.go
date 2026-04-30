package portlookup

import "flag"

// Config controls Lookup behaviour.
type Config struct {
	// CacheEnabled controls whether resolved results are cached.
	CacheEnabled bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		CacheEnabled: true,
	}
}

// RegisterFlags registers portlookup flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.CacheEnabled, "lookup-cache", cfg.CacheEnabled,
		"cache resolved service names for the lifetime of the scan")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if !src.CacheEnabled {
		dst.CacheEnabled = false
	}
}
