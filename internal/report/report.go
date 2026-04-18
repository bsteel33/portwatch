package report

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a snapshot and diff for rendering.
type Report struct {
	Timestamp time.Time
	Snapshot  *snapshot.Snapshot
	Diff      *snapshot.Diff
	Format    Format
	Writer    io.Writer
}

// New creates a new Report with default settings.
func New(snap *snapshot.Snapshot, diff *snapshot.Diff) *Report {
	return &Report{
		Timestamp: time.Now(),
		Snapshot:  snap,
		Diff:      diff,
		Format:    FormatText,
		Writer:    os.Stdout,
	}
}

// Render writes the report to the configured writer.
func (r *Report) Render() error {
	switch r.Format {
	case FormatJSON:
		return r.renderJSON()
	default:
		return r.renderText()
	}
}

func (r *Report) renderText() error {
	fmt.Fprintf(r.Writer, "=== portwatch report [%s] ===\n", r.Timestamp.Format(time.RFC3339))
	if r.Diff != nil && (len(r.Diff.Added) > 0 || len(r.Diff.Removed) > 0) {
		snapshot.PrintDiff(r.Writer, r.Diff)
	} else {
		fmt.Fprintln(r.Writer, "No changes detected.")
	}
	if r.Snapshot != nil {
		fmt.Fprintln(r.Writer, "--- current open ports ---")
		snapshot.PrintSnapshot(r.Writer, r.Snapshot)
	}
	return nil
}

func (r *Report) renderJSON() error {
	enc := jsonEncoder(r.Writer)
	return enc.Encode(map[string]interface{}{
		"timestamp": r.Timestamp.Format(time.RFC3339),
		"added":     r.Diff.Added,
		"removed":   r.Diff.Removed,
		"ports":     r.Snapshot.Ports,
	})
}
