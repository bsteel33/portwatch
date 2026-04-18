package baseline

import "flag"

// Config holds configuration for the baseline manager.
type Config struct {
	Path string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path: ".portwatch_baseline.json",
	}
}

// RegisterFlags registers baseline-related CLI flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "baseline", cfg.Path, "path to baseline file")
}

// ApplyFlags overwrites cfg fields with non-zero values from src.
func ApplyFlags(cfg *Config, src Config) {
	if src.Path != "" {
		cfg.Path = src.Path
	}
}
