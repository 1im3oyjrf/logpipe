// Package merge provides a processor that merges fields from a static map
// into each log entry, optionally overwriting existing keys.
package merge

import (
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls which fields are merged and whether existing keys are
// overwritten.
type Config struct {
	// Fields is the set of key/value pairs to merge into every entry.
	Fields map[string]string
	// Overwrite controls whether an existing field is replaced.
	Overwrite bool
}

// Merger merges a fixed set of fields into log entries.
type Merger struct {
	cfg Config
}

// New returns a Merger configured with cfg.
func New(cfg Config) *Merger {
	normalised := make(map[string]string, len(cfg.Fields))
	for k, v := range cfg.Fields {
		normalised[strings.ToLower(k)] = v
	}
	return &Merger{cfg: Config{Fields: normalised, Overwrite: cfg.Overwrite}}
}

// Apply returns a shallow copy of entry with the configured fields merged in.
func (m *Merger) Apply(entry map[string]any) map[string]any {
	if len(m.cfg.Fields) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry)+len(m.cfg.Fields))
	for k, v := range entry {
		out[k] = v
	}
	for k, v := range m.cfg.Fields {
		existing := parser.HasField(out, k)
		if existing && !m.cfg.Overwrite {
			continue
		}
		out[k] = v
	}
	return out
}
