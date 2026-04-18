// Package stamp rewrites or injects a timestamp field on each log entry.
package stamp

import (
	"time"

	"github.com/yourorg/logpipe/internal/parser"
)

// Config controls the stamping behaviour.
type Config struct {
	// Field is the key written to the entry (default: "timestamp").
	Field string
	// Format is a Go time layout string (default: time.RFC3339).
	Format string
	// Overwrite replaces an existing field when true.
	Overwrite bool
}

// Stamper rewrites or injects a timestamp field.
type Stamper struct {
	field     string
	format    string
	overwrite bool
	now       func() time.Time
}

// New returns a Stamper configured by cfg.
func New(cfg Config) *Stamper {
	field := cfg.Field
	if field == "" {
		field = "timestamp"
	}
	fmt := cfg.Format
	if fmt == "" {
		fmt = time.RFC3339
	}
	return &Stamper{
		field:     field,
		format:    fmt,
		overwrite: cfg.Overwrite,
		now:       time.Now,
	}
}

// Apply returns a copy of entry with the timestamp field set.
func (s *Stamper) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}

	key := parser.CanonicalKey(s.field, out)
	if key == "" {
		key = s.field
	}

	if _, exists := out[key]; exists && !s.overwrite {
		return out
	}

	out[key] = s.now().Format(s.format)
	return out
}
