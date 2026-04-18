package history

import "flag"

const defaultHistoryPath = "/var/lib/portwatch/history.json"

// Config holds configuration for the history module.
type Config struct {
	Path    string
	Enabled bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Path:    defaultHistoryPath,
		Enabled: true,
	}
}

// ApplyFlags overrides Config fields from parsed CLI flags.
func ApplyFlags(c *Config, fs *flag.FlagSet) {
	if f := fs.Lookup("history-path"); f != nil && f.Value.String() != "" {
		c.Path = f.Value.String()
	}
	if f := fs.Lookup("no-history"); f != nil && f.Value.String() == "true" {
		c.Enabled = false
	}
}

// RegisterFlags registers history-related flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet) {
	fs.String("history-path", "", "path to history file (default: "+defaultHistoryPath+")")
	fs.Bool("no-history", false, "disable history recording")
}
