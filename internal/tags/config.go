package tags

import "flag"

// Config holds configuration for the tags module.
type Config struct {
	Path string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path: "tags.json",
	}
}

// RegisterFlags registers tag-related CLI flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "tags", cfg.Path, "path to port tags JSON file")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
}
