// Package clip trims leading and trailing whitespace from string fields
// in a log entry, optionally targeting specific fields only.
package clip

import (
	"strings"

	"logpipe/internal/parser"
)

// Config controls which fields are trimmed.
type Config struct {
	// Fields lists the field names to trim. When empty all string fields are trimmed.
	Fields []string
}

// Clipper trims whitespace from log entry fields.
type Clipper struct {
	fields map[string]struct{}
	all    bool
}

// New returns a Clipper configured by cfg.
func New(cfg Config) *Clipper {
	c := &Clipper{fields: make(map[string]struct{})}
	if len(cfg.Fields) == 0 {
		c.all = true
		return c
	}
	for _, f := range cfg.Fields {
		c.fields[strings.ToLower(f)] = struct{}{}
	}
	return c
}

// Apply returns a copy of entry with whitespace trimmed from targeted string fields.
func (c *Clipper) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, k := range parser.Keys(entry) {
		s, ok := parser.GetString(entry, k)
		if !ok {
			continue
		}
		if c.all {
			out[k] = strings.TrimSpace(s)
			continue
		}
		if _, targeted := c.fields[strings.ToLower(k)]; targeted {
			out[k] = strings.TrimSpace(s)
		}
	}
	return out
}
