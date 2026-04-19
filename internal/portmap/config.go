package portmap

import "flag"

// Config holds portmap configuration.
type Config struct {
	Path string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path: "portmap.json",
	}
}

// RegisterFlags registers portmap flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "portmap-file", cfg.Path, "path to custom port name map file")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
}
