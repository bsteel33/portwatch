package audit

import "flag"

// Config holds audit logger configuration.
type Config struct {
	Path    string
	Enabled bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:    "portwatch-audit.log",
		Enabled: true,
	}
}

// RegisterFlags registers audit flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.Path, "audit-log", cfg.Path, "path to audit log file")
	fs.BoolVar(&cfg.Enabled, "audit", cfg.Enabled, "enable audit logging")
}

// ApplyFlags copies non-zero values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Path != "" {
		dst.Path = src.Path
	}
	dst.Enabled = src.Enabled
}

// Validate checks that the Config is valid, returning an error if not.
// A Config is invalid if audit logging is enabled but no path is specified.
func (c Config) Validate() error {
	if c.Enabled && c.Path == "" {
		return fmt.Errorf("audit: log path must not be empty when audit logging is enabled")
	}
	return nil
}
