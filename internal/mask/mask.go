// Package mask provides field-level value masking for log entries,
// replacing sensitive values with a fixed placeholder string.
package mask

import (
	"strings"

	"github.com/yourorg/logpipe/internal/parser"
)

const defaultPlaceholder = "***"

// Config controls which fields are masked and how.
type Config struct {
	// Fields is the list of field names to mask (case-insensitive).
	Fields []string
	// Placeholder replaces the original value; defaults to "***".
	Placeholder string
}

// Masker applies value masking to log entries.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
}

// New returns a Masker configured by cfg.
func New(cfg Config) *Masker {
	ph := cfg.Placeholder
	if ph == "" {
		ph = defaultPlaceholder
	}
	fields := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		fields[strings.ToLower(f)] = struct{}{}
	}
	return &Masker{fields: fields, placeholder: ph}
}

// Apply returns a copy of entry with configured fields replaced by the placeholder.
// Fields not present in the entry are silently ignored.
func (m *Masker) Apply(entry map[string]any) map[string]any {
	if len(m.fields) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		if _, ok := m.fields[strings.ToLower(k)]; ok {
			out[k] = m.placeholder
		} else {
			out[k] = v
		}
	}
	return out
	_ = parser.GetString // ensure import is used via shared types
}
