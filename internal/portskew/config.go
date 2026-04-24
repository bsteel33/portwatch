package portskew

import "flag"

// RegisterFlags registers portskew flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Float64Var(&cfg.Threshold, "skew-threshold", cfg.Threshold,
		"z-score threshold for flagging port-count skew")
	fs.IntVar(&cfg.MinSamples, "skew-min-samples", cfg.MinSamples,
		"minimum samples required before skew detection activates")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.Threshold != 0 {
		dst.Threshold = src.Threshold
	}
	if src.MinSamples != 0 {
		dst.MinSamples = src.MinSamples
	}
}
