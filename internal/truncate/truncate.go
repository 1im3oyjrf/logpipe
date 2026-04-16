// Package truncate provides field-level value truncation for log entries.
package truncate

import (
	"strings"

	"github.com/your-org/logpipe/internal/parser"
)

// Config controls which fields are truncated and to what length.
type Config struct {
	// MaxLen is the default maximum byte length applied to string values.
	MaxLen int
	// Fields lists specific field names to truncate (case-insensitive).
	// When empty every string field is subject to MaxLen.
	Fields []string
}

// Truncator applies truncation rules to log entries.
type Truncator struct {
	cfg    Config
	fields map[string]struct{}
}

// New creates a Truncator from cfg. A MaxLen <= 0 defaults to 256.
func New(cfg Config) *Truncator {
	if cfg.MaxLen <= 0 {
		cfg.MaxLen = 256
	}
	fields := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		fields[strings.ToLower(f)] = struct{}{}
	}
	return &Truncator{cfg: cfg, fields: fields}
}

// Apply returns a new entry with string values truncated according to cfg.
// The original entry is never mutated.
func (t *Truncator) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, k := range parser.Keys(entry) {
		if !t.shouldTruncate(k) {
			continue
		}
		s, ok := parser.GetString(entry, k)
		if !ok {
			continue
		}
		if len(s) > t.cfg.MaxLen {
			out[k] = s[:t.cfg.MaxLen]
		}
	}
	return out
}

func (t *Truncator) shouldTruncate(field string) bool {
	if len(t.fields) == 0 {
		return true
	}
	_, ok := t.fields[strings.ToLower(field)]
	return ok
}
