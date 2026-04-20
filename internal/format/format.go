// Package format provides a field value formatter that applies
// sprintf-style templates to produce a new field from existing entry data.
package format

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls how the formatter rewrites entries.
type Config struct {
	// Target is the field name to write the formatted value into.
	Target string

	// Template is a Go fmt-style template where {field} placeholders are
	// replaced with the corresponding field value from the log entry.
	Template string

	// Overwrite controls whether an existing target field is replaced.
	Overwrite bool
}

// Formatter rewrites a log entry by rendering a template into a target field.
type Formatter struct {
	cfg Config
}

// New returns a Formatter for the given Config.
// An error is returned when Target or Template is empty.
func New(cfg Config) (*Formatter, error) {
	if strings.TrimSpace(cfg.Target) == "" {
		return nil, fmt.Errorf("format: target field must not be empty")
	}
	if strings.TrimSpace(cfg.Template) == "" {
		return nil, fmt.Errorf("format: template must not be empty")
	}
	return &Formatter{cfg: cfg}, nil
}

// Apply renders the template against entry and returns a new entry with the
// target field set. The original entry is never mutated.
func (f *Formatter) Apply(entry map[string]any) map[string]any {
	targetKey := strings.ToLower(f.cfg.Target)

	// Respect Overwrite flag.
	if !f.cfg.Overwrite {
		for k := range entry {
			if strings.ToLower(k) == targetKey {
				return entry
			}
		}
	}

	result := render(f.cfg.Template, entry)

	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	out[f.cfg.Target] = result
	return out
}

// render replaces {field} tokens in tmpl with values from entry.
func render(tmpl string, entry map[string]any) string {
	result := tmpl
	for k, v := range entry {
		placeholder := "{" + k + "}"
		result = strings.ReplaceAll(result, placeholder, parser.GetString(entry, k))
		_ = v
	}
	return result
}
