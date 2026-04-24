package portresolve

import "flag"

// Config holds configuration for the Resolver.
type Config struct {
	// FallbackPrefix is used when no name is found. Defaults to "port".
	FallbackPrefix string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		FallbackPrefix: "port",
	}
}

// RegisterFlags registers portresolve flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.FallbackPrefix, "resolve.fallback-prefix", cfg.FallbackPrefix,
		"prefix used for unknown port names (e.g. \"port\" → \"port-8080\")")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.FallbackPrefix != "" {
		dst.FallbackPrefix = src.FallbackPrefix
	}
}
