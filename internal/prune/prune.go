// Package prune removes fields from log entries based on configuration.
package prune

import "github.com/logpipe/logpipe/internal/parser"

// Config controls which fields are removed from each entry.
type Config struct {
	// Fields lists exact field names to remove (case-insensitive).
	Fields []string
}

// Pruner removes configured fields from log entries.
type Pruner struct {
	fields map[string]struct{}
}

// New returns a Pruner that drops the fields listed in cfg.
func New(cfg Config) *Pruner {
	m := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		m[strings.ToLower(f)] = struct{}{}
	}
	return &Pruner{fields: m}
}

// Apply returns a shallow copy of entry with configured fields removed.
// The original entry is never mutated.
func (p *Pruner) Apply(entry map[string]any) map[string]any {
	if len(p.fields) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry))
	for _, k := range parser.Keys(entry) {
		if _, drop := p.fields[strings.ToLower(k)]; !drop {
			out[k] = entry[k]
		}
	}
	return out
}
