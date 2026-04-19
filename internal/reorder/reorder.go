// Package reorder provides a stage that reorders fields in a log entry
// so that a configured set of keys appear first in the output map.
package reorder

import (
	"strings"

	"logpipe/internal/reader"
)

// Config controls which fields are promoted to the front.
type Config struct {
	// Fields lists the keys that should appear first, in order.
	Fields []string

	// CaseInsensitive controls whether field matching ignores case.
	CaseInsensitive bool
}

// Reorder promotes a configured set of fields to the front of each entry's
// field map and passes the result downstream.
type Reorder struct {
	cfg Config
}

// New returns a Reorder stage configured with cfg.
func New(cfg Config) *Reorder {
	return &Reorder{cfg: cfg}
}

// Apply returns a new entry whose fields begin with the configured keys
// (in order) followed by all remaining fields in their original order.
func (r *Reorder) Apply(e reader.Entry) reader.Entry {
	if len(r.cfg.Fields) == 0 {
		return e
	}

	out := make(map[string]any, len(e.Fields))

	// Track which original keys have been placed already.
	placed := make(map[string]bool, len(r.cfg.Fields))

	for _, want := range r.cfg.Fields {
		for k, v := range e.Fields {
			if r.matches(k, want) {
				out[k] = v
				placed[k] = true
				break
			}
		}
	}

	for k, v := range e.Fields {
		if !placed[k] {
			out[k] = v
		}
	}

	return reader.Entry{
		Level:   e.Level,
		Message: e.Message,
		Fields:  out,
		Source:  e.Source,
	}
}

func (r *Reorder) matches(key, want string) bool {
	if r.cfg.CaseInsensitive {
		return strings.EqualFold(key, want)
	}
	return key == want
}
