package portcheck

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config holds configuration for the Checker.
type Config struct {
	// Raw rule strings in the form "port/proto:open" or "port/proto:closed".
	Rules []string
}

// DefaultConfig returns a Config with no conditions.
func DefaultConfig() Config {
	return Config{}
}

// RegisterFlags registers CLI flags on fs and returns a pointer to the Config
// that will be populated after parsing.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Func("check", `health check rule: port/proto:open|closed (e.g. 80/tcp:open)`, func(s string) error {
		_, err := parseRule(s)
		if err != nil {
			return err
		}
		cfg.Rules = append(cfg.Rules, s)
		return nil
	})
}

// Build converts the Config into a slice of Conditions.
func Build(cfg Config) ([]Condition, error) {
	conds := make([]Condition, 0, len(cfg.Rules))
	for _, r := range cfg.Rules {
		c, err := parseRule(r)
		if err != nil {
			return nil, err
		}
		conds = append(conds, c)
	}
	return conds, nil
}

// parseRule parses a rule string of the form "port/proto:open" or "port/proto:closed".
func parseRule(s string) (Condition, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Condition{}, fmt.Errorf("portcheck: invalid rule %q: expected port/proto:open|closed", s)
	}
	portProto := strings.SplitN(parts[0], "/", 2)
	if len(portProto) != 2 {
		return Condition{}, fmt.Errorf("portcheck: invalid rule %q: expected port/proto", s)
	}
	port, err := strconv.Atoi(portProto[0])
	if err != nil || port < 1 || port > 65535 {
		return Condition{}, fmt.Errorf("portcheck: invalid port %q in rule %q", portProto[0], s)
	}
	proto := strings.ToLower(portProto[1])
	if proto != "tcp" && proto != "udp" {
		return Condition{}, fmt.Errorf("portcheck: invalid proto %q in rule %q", proto, s)
	}
	switch strings.ToLower(parts[1]) {
	case "open":
		return Condition{Port: port, Proto: proto, MustBeOpen: true}, nil
	case "closed":
		return Condition{Port: port, Proto: proto, MustBeOpen: false}, nil
	default:
		return Condition{}, fmt.Errorf("portcheck: invalid state %q in rule %q: must be open or closed", parts[1], s)
	}
}
