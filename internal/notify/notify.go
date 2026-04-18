package notify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/snapshot"
)

// Channel represents a notification delivery method.
type Channel string

const (
	ChannelStdout Channel = "stdout"
	ChannelExec   Channel = "exec"
)

// Config holds configuration for a notifier.
type Config struct {
	Channel Channel
	// ExecCmd is the command to run when Channel == ChannelExec.
	// The diff summary is passed as the first argument.
	ExecCmd string
}

// Notifier sends notifications when port changes are detected.
type Notifier struct {
	cfg Config
	out io.Writer
}

// New creates a new Notifier.
func New(cfg Config) *Notifier {
	return &Notifier{cfg: cfg, out: os.Stdout}
}

// Notify dispatches a notification for the given diff.
func (n *Notifier) Notify(diff snapshot.Diff) error {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return nil
	}
	msg := buildMessage(diff)
	switch n.cfg.Channel {
	case ChannelExec:
		return runCmd(n.cfg.ExecCmd, msg)
	default:
		fmt.Fprintln(n.out, msg)
		return nil
	}
}

func buildMessage(diff snapshot.Diff) string {
	var sb strings.Builder
	sb.WriteString("[portwatch] Port changes detected:\n")
	for _, p := range diff.Added {
		fmt.Fprintf(&sb, "  + %d/%s (%s)\n", p.Port, p.Proto, p.Service)
	}
	for _, p := range diff.Removed {
		fmt.Fprintf(&sb, "  - %d/%s (%s)\n", p.Port, p.Proto, p.Service)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCmd(cmd, msg string) error {
	if cmd == "" {
		return fmt.Errorf("notify: exec command is empty")
	}
	c := exec.Command(cmd, msg)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
