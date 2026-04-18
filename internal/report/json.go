package report

import (
	"encoding/json"
	"io"
)

// jsonEncoder returns a configured JSON encoder for report output.
func jsonEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc
}
