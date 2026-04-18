package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelCrit  Level = "CRIT"
)

// Alert holds a single alert event.
type Alert struct {
	Time    time.Time
	Level   Level
	Message string
}

// Notifier writes alerts to an output destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to w. Pass nil to use os.Stdout.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify converts a snapshot diff into alerts and writes them.
func (n *Notifier) Notify(diff snapshot.Diff) []Alert {
	var alerts []Alert

	for _, p := range diff.Added {
		a := Alert{
			Time:    time.Now(),
			Level:   LevelWarn,
			Message: fmt.Sprintf("port opened: %d/%s (%s)", p.Port, p.Proto, p.Service),
		}
		alerts = append(alerts, a)
		n.write(a)
	}

	for _, p := range diff.Removed {
		a := Alert{
			Time:    time.Now(),
			Level:   LevelInfo,
			Message: fmt.Sprintf("port closed: %d/%s (%s)", p.Port, p.Proto, p.Service),
		}
		alerts = append(alerts, a)
		n.write(a)
	}

	return alerts
}

func (n *Notifier) write(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s %s\n", a.Level, a.Time.Format(time.RFC3339), a.Message)
}
