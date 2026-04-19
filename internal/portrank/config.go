package portrank

import "flag"

// Config controls ranking behaviour.
type Config struct {
	// Weights maps port numbers to additional risk scores.
	Weights map[int]int
}

// DefaultConfig returns sensible defaults with common high-risk ports pre-weighted.
func DefaultConfig() Config {
	return Config{
		Weights: map[int]int{
			21:   70, // FTP
			22:   40, // SSH
			23:   90, // Telnet
			25:   60, // SMTP
			445:  80, // SMB
			3389: 85, // RDP
			5900: 75, // VNC
		},
	}
}

// RegisterFlags registers portrank flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet) {
	// Future: allow users to supply custom weight overrides via CLI.
	_ = fs
}

// ApplyFlags merges parsed flag values into cfg.
func ApplyFlags(cfg *Config, fs *flag.FlagSet) {
	_ = fs
}
