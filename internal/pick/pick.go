// Package pick provides a transformer that retains only a configured set of
// fields from each log entry, discarding everything else.
package pick

import (
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls which fields are kept.
type Config struct {
	// Fields lists the field names to retain (case-insensitive).
	Fields []string
}

// Picker retains only the configured fields from each entry.
type Picker struct {
	fields map[string]struct{}
}

// New returns a Picker configured to keep only the named fields.
// If no fields are specified, entries pass through unchanged.
func New(cfg Config) *Picker {
	set := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		set[strings.ToLower(f)] = struct{}{}
	}
	return &Picker{fields: set}
}

// Apply returns a new entry containing only the configured fields.
// If the Picker has no configured fields the original entry is returned as-is.
func (p *Picker) Apply(entry map[string]any) map[string]any {
	if len(p.fields) == 0 {
		return entry
	}
	out := make(map[string]any, len(p.fields))
	for _, k := range parser.Keys(entry) {
		if _, ok := p.fields[strings.ToLower(k)]; ok {
			out[k] = entry[k]
		}
	}
	return out
}
