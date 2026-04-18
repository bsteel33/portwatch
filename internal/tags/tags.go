// Package tags provides port tagging — associating human-readable labels
// with specific ports or port ranges for richer reporting output.
package tags

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// Tag associates a label with a port number.
type Tag struct {
	Port  int    `json:"port"`
	Proto string `json:"proto"` // tcp or udp
	Label string `json:"label"`
}

// Tagger holds a set of tags and resolves labels for ports.
type Tagger struct {
	tags []Tag
}

// New returns a Tagger loaded from the given file path.
// If path is empty, an empty Tagger is returned.
func New(path string) (*Tagger, error) {
	if path == "" {
		return &Tagger{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Tagger{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var tags []Tag
	if err := json.NewDecoder(f).Decode(&tags); err != nil {
		return nil, err
	}
	return &Tagger{tags: tags}, nil
}

// Resolve returns the label for the given port and proto, or an empty string.
func (t *Tagger) Resolve(port int, proto string) string {
	proto = strings.ToLower(proto)
	for _, tag := range t.tags {
		if tag.Port == port && strings.ToLower(tag.Proto) == proto {
			return tag.Label
		}
	}
	return ""
}

// Add appends a tag to the in-memory set.
func (t *Tagger) Add(port int, proto, label string) {
	t.tags = append(t.tags, Tag{Port: port, Proto: strings.ToLower(proto), Label: label})
}

// Save writes the current tag set to path as JSON.
func (t *Tagger) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(t.tags)
}

// Key returns a canonical key string for a port+proto pair.
func Key(port int, proto string) string {
	return strconv.Itoa(port) + "/" + strings.ToLower(proto)
}
