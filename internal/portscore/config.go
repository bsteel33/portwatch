package portscore

import "flag"

// Config controls how individual risk factors are weighted.
type Config struct {
	WellKnownWeight int
	TCPWeight       int
	// HighRiskPorts maps "port/proto" strings to additional weight.
	HighRiskPorts map[string]int
}

// DefaultConfig returns a sensible default Config.
func DefaultConfig() Config {
	return Config{
		WellKnownWeight: 20,
		TCPWeight:       10,
		HighRiskPorts: map[string]int{
			"23/tcp":  50, // telnet
			"21/tcp":  40, // ftp
			"3389/tcp": 45, // rdp
			"445/tcp": 45, // smb
		},
	}
}

// RegisterFlags registers portscore flags on fs.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.WellKnownWeight, "score.well-known-weight", cfg.WellKnownWeight,
		"risk weight added for ports below 1024")
	fs.IntVar(&cfg.TCPWeight, "score.tcp-weight", cfg.TCPWeight,
		"risk weight added for TCP ports")
}

// ApplyFlags copies non-zero flag values from src into dst.
func ApplyFlags(dst *Config, src Config) {
	if src.WellKnownWeight != 0 {
		dst.WellKnownWeight = src.WellKnownWeight
	}
	if src.TCPWeight != 0 {
		dst.TCPWeight = src.TCPWeight
	}
}
