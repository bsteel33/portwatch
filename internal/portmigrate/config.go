package portmigrate

import "flag"

// Config holds configuration for the migration module.
type Config struct {
	// SnapshotPath is the file to read and migrate.
	SnapshotPath string
	// DryRun reports what would be migrated without writing changes.
	DryRun bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		SnapshotPath: "portwatch.snap",
		DryRun:       false,
	}
}

// RegisterFlags binds Config fields to the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.SnapshotPath, "migrate.snapshot", cfg.SnapshotPath,
		"path to snapshot file to migrate")
	fs.BoolVar(&cfg.DryRun, "migrate.dry-run", cfg.DryRun,
		"report migration steps without writing output")
}

// ApplyFlags copies non-zero override values onto cfg.
func ApplyFlags(cfg *Config, snapshot string, dryRun bool) {
	if snapshot != "" {
		cfg.SnapshotPath = snapshot
	}
	if dryRun {
		cfg.DryRun = true
	}
}
