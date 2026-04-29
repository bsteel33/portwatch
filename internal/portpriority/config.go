package portpriority

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config holds configuration for the Prioritizer.
type Config struct {
	// Rules is a list of "port/proto=level" strings, e.g. "22/tcp=critical".
	Rules []string
	// DefaultLevel is the fallback level name when no rule matches.
	DefaultLevel string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultLevel: "low",
	}
}

// RegisterFlags registers CLI flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&cfg.DefaultLevel, "priority.default", cfg.DefaultLevel, "default priority level (low|medium|high|critical)")
	fs.Func("priority.rule", "port priority rule: port/proto=level (repeatable)", func(s string) error {
		cfg.Rules = append(cfg.Rules, s)
		return nil
	})
}

// Build constructs a Prioritizer from the Config.
func Build(cfg Config) (*Prioritizer, error) {
	defLevel, err := parseLevel(cfg.DefaultLevel)
	if err != nil {
		return nil, fmt.Errorf("priority.default: %w", err)
	}
	var rules []Rule
	for _, raw := range cfg.Rules {
		r, err := ParseRule(raw)
		if err != nil {
			return nil, fmt.Errorf("priority.rule %q: %w", raw, err)
		}
		rules = append(rules, r)
	}
	return New(rules, defLevel), nil
}

// ParseRule parses a rule string of the form "port/proto=level" or "port=level".
func ParseRule(s string) (Rule, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("expected port[/proto]=level, got %q", s)
	}
	lvl, err := parseLevel(parts[1])
	if err != nil {
		return Rule{}, err
	}
	pp := strings.SplitN(parts[0], "/", 2)
	port, err := strconv.Atoi(pp[0])
	if err != nil {
		return Rule{}, fmt.Errorf("invalid port %q", pp[0])
	}
	proto := ""
	if len(pp) == 2 {
		proto = strings.ToLower(pp[1])
	}
	return Rule{Port: port, Proto: proto, Level: lvl}, nil
}

func parseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "low":
		return Low, nil
	case "medium":
		return Medium, nil
	case "high":
		return High, nil
	case "critical":
		return Critical, nil
	default:
		return Low, fmt.Errorf("unknown level %q", s)
	}
}
