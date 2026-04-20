package portclassify

import "flag"

// Rule defines a single classification rule.
type Rule struct {
	Port   int
	Proto  string
	Class  Class
	Reason string
}

// Config holds classifier configuration.
type Config struct {
	Rules []Rule
}

// DefaultConfig returns a Config pre-loaded with common risk rules.
func DefaultConfig() Config {
	return Config{
		Rules: []Rule{
			{Port: 22, Proto: "tcp", Class: ClassMonitor, Reason: "SSH – ensure key-only auth"},
			{Port: 23, Proto: "tcp", Class: ClassDangerous, Reason: "Telnet – plaintext protocol"},
			{Port: 80, Proto: "tcp", Class: ClassSafe, Reason: "HTTP"},
			{Port: 443, Proto: "tcp", Class: ClassSafe, Reason: "HTTPS"},
			{Port: 3306, Proto: "tcp", Class: ClassSuspicious, Reason: "MySQL exposed"},
			{Port: 5432, Proto: "tcp", Class: ClassSuspicious, Reason: "PostgreSQL exposed"},
			{Port: 6379, Proto: "tcp", Class: ClassDangerous, Reason: "Redis – often unauthenticated"},
			{Port: 27017, Proto: "tcp", Class: ClassDangerous, Reason: "MongoDB – often unauthenticated"},
			{Port: 21, Proto: "tcp", Class: ClassSuspicious, Reason: "FTP – plaintext"},
			{Port: 25, Proto: "tcp", Class: ClassMonitor, Reason: "SMTP"},
			{Port: 53, Proto: "udp", Class: ClassMonitor, Reason: "DNS"},
			{Port: 8080, Proto: "tcp", Class: ClassMonitor, Reason: "HTTP-alt"},
			{Port: 8443, Proto: "tcp", Class: ClassMonitor, Reason: "HTTPS-alt"},
		},
	}
}

// RegisterFlags registers classifier-related flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, _ *Config) {
	// Future: --classify-rules-file path
	_ = fs
}

// ApplyFlags copies flag values into cfg (no-op until custom rule files are supported).
func ApplyFlags(_ *flag.FlagSet, _ *Config) {}
