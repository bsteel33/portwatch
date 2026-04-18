package notify

import (
	"fmt"
	"strings"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Channel: ChannelStdout,
	}
}

// ParseChannel converts a string to a Channel, returning an error for unknown values.
func ParseChannel(s string) (Channel, error) {
	switch strings.ToLower(s) {
	case "stdout", "":
		return ChannelStdout, nil
	case "exec":
		return ChannelExec, nil
	default:
		return "", fmt.Errorf("notify: unknown channel %q (valid: stdout, exec)", s)
	}
}

// ApplyFlags merges flag-provided values into a Config.
func ApplyFlags(cfg *Config, channel, execCmd string) error {
	if channel != "" {
		ch, err := ParseChannel(channel)
		if err != nil {
			return err
		}
		cfg.Channel = ch
	}
	if execCmd != "" {
		cfg.ExecCmd = execCmd
	}
	return nil
}
