// Package promote lifts a nested field value to the top level of a log entry.
// If the source field contains a map, its keys are merged into the entry.
// If it contains a scalar value, it is assigned to a configurable target key.
package promote

import (
	"fmt"
	"strings"

	"logpipe/internal/parser"
)

// Config controls which field is promoted and how scalar values are handled.
type Config struct {
	// Field is the dot-separated path to the field whose value should be promoted.
	Field string

	// Target is the top-level key used when the promoted value is a scalar.
	// Defaults to the leaf segment of Field.
	Target string

	// DropSource removes the original field after promotion. Defaults to true.
	DropSource bool

	// CaseInsensitive enables case-insensitive field matching.
	CaseInsensitive bool
}

// Promoter applies field promotion to log entries.
type Promoter struct {
	cfg Config
}

// New returns a Promoter for the given Config, or an error if the config is invalid.
func New(cfg Config) (*Promoter, error) {
	cfg.Field = strings.TrimSpace(cfg.Field)
	if cfg.Field == "" {
		return nil, fmt.Errorf("promote: field must not be empty")
	}
	if cfg.Target == "" {
		parts := strings.Split(cfg.Field, ".")
		cfg.Target = parts[len(parts)-1]
	}
	return &Promoter{cfg: cfg}, nil
}

// Apply promotes the configured field within entry and returns the modified copy.
func (p *Promoter) Apply(entry map[string]any) map[string]any {
	key := p.resolveKey(entry)
	if key == "" {
		return entry
	}

	out := shallowCopy(entry)
	val := out[key]

	if p.cfg.DropSource {
		delete(out, key)
	}

	switch v := val.(type) {
	case map[string]any:
		for k, mv := range v {
			out[k] = mv
		}
	default:
		out[p.cfg.Target] = val
	}

	return out
}

func (p *Promoter) resolveKey(entry map[string]any) string {
	if !p.cfg.CaseInsensitive {
		if _, ok := entry[p.cfg.Field]; ok {
			return p.cfg.Field
		}
		return ""
	}
	return parser.FindKey(entry, p.cfg.Field)
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
