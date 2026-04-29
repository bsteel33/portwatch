package portrelay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// WriterDestination writes JSON-encoded snapshots to an io.Writer.
type WriterDestination struct {
	name string
	w    io.Writer
}

// NewWriterDestination creates a destination that writes to w.
func NewWriterDestination(name string, w io.Writer) *WriterDestination {
	return &WriterDestination{name: name, w: w}
}

func (d *WriterDestination) Name() string { return d.name }

func (d *WriterDestination) Send(snap *snapshot.Snapshot) error {
	return json.NewEncoder(d.w).Encode(snap)
}

// HTTPDestination POSTs JSON-encoded snapshots to a URL.
type HTTPDestination struct {
	name   string
	url    string
	client *http.Client
}

// NewHTTPDestination creates a destination that POSTs to url.
func NewHTTPDestination(name, url string, timeout time.Duration) *HTTPDestination {
	return &HTTPDestination{
		name:   name,
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

func (d *HTTPDestination) Name() string { return d.name }

func (d *HTTPDestination) Send(snap *snapshot.Snapshot) error {
	body, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	resp, err := d.client.Post(d.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}
