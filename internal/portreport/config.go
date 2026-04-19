package portreport

import "flag"

// Config controls report generation behaviour.
type Config struct {
	ShowUnknown bool
	SortByScore bool
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		ShowUnknown: true,
		SortByScore: false,
	}
}

// RegisterFlags registers CLI flags into the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.ShowUnknown, "report.show-unknown", cfg.ShowUnknown, "include ports with unknown service names in report")
	fs.BoolVar(&cfg.SortByScore, "report.sort-score", cfg.SortByScore, "sort report entries by risk score descending")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	dst.ShowUnknown = src.ShowUnknown
	dst.SortByScore = src.SortByScore
}
