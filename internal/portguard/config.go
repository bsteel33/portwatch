package portguard

import "flag"

// Config holds configuration for the Guard.
type Config struct {
	// Allowlist is a list of "port/proto" strings that are explicitly allowed.
	Allowlist []string

	// Denylist is a list of "port/proto" strings that are explicitly denied.
	Denylist []string

	// Default is the action taken when a port matches neither list.
	// Defaults to ActionAllow.
	Default Action
}

// DefaultConfig returns a permissive guard configuration.
func DefaultConfig() Config {
	return Config{
		Default: ActionAllow,
	}
}

// RegisterFlags registers portguard flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Func("guard-allow", "allow a port (format: port/proto, repeatable)", func(s string) error {
		cfg.Allowlist = append(cfg.Allowlist, s)
		return nil
	})
	fs.Func("guard-deny", "deny a port (format: port/proto, repeatable)", func(s string) error {
		cfg.Denylist = append(cfg.Denylist, s)
		return nil
	})
	fs.Func("guard-default", "default action when no rule matches: allow|deny (default: allow)", func(s string) error {
		switch Action(s) {
		case ActionAllow, ActionDeny:
			cfg.Default = Action(s)
			return nil
		default:
			return fmt.Errorf("invalid guard-default %q: must be allow or deny", s)
		}
	})
}

// ApplyFlags overwrites cfg fields only when non-zero values are provided.
func ApplyFlags(cfg *Config, allowlist, denylist []string, defaultAction string) {
	if len(allowlist) > 0 {
		cfg.Allowlist = allowlist
	}
	if len(denylist) > 0 {
		cfg.Denylist = denylist
	}
	if defaultAction != "" {
		cfg.Default = Action(defaultAction)
	}
}
