package portshield

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config holds configuration for the shield.
type Config struct {
	// AllowPorts is a comma-separated list of port/proto pairs to allow, e.g. "80/tcp,443/tcp".
	AllowPorts string
	// BlockPorts is a comma-separated list of port/proto pairs to block.
	BlockPorts string
	// DefaultAction is the fallback action ("allow" or "block").
	DefaultAction string
}

// DefaultConfig returns a Config that allows everything by default.
func DefaultConfig() Config {
	return Config{DefaultAction: "allow"}
}

// RegisterFlags registers shield-related flags on the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.AllowPorts, "shield-allow", cfg.AllowPorts, "comma-separated port/proto pairs to explicitly allow (e.g. 80/tcp)")
	fs.StringVar(&cfg.BlockPorts, "shield-block", cfg.BlockPorts, "comma-separated port/proto pairs to block")
	fs.StringVar(&cfg.DefaultAction, "shield-default", cfg.DefaultAction, "default shield action: allow or block")
}

// Build constructs a Shield from cfg.
func Build(cfg Config) (*Shield, error) {
	defaultAction := Allow
	if strings.ToLower(cfg.DefaultAction) == "block" {
		defaultAction = Block
	}
	s := New(defaultAction)
	if err := applyRules(s, cfg.AllowPorts, Allow); err != nil {
		return nil, err
	}
	if err := applyRules(s, cfg.BlockPorts, Block); err != nil {
		return nil, err
	}
	return s, nil
}

func applyRules(s *Shield, raw string, action Action) error {
	if raw == "" {
		return nil
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		port, proto, err := parsePortProto(part)
		if err != nil {
			return err
		}
		s.Add(port, proto, action)
	}
	return nil
}

func parsePortProto(s string) (int, string, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("portshield: invalid rule %q, expected port/proto", s)
	}
	port, err := strconv.Atoi(parts[0])
	if err != nil || port < 1 || port > 65535 {
		return 0, "", fmt.Errorf("portshield: invalid port in rule %q", s)
	}
	return port, strings.ToLower(parts[1]), nil
}
