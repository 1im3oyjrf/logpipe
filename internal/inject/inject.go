package inject

import (
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config holds the configuration for the field injector.
type Config struct {
	// Fields is a map of field name to static value to inject into every entry.
	Fields map[string]string
	// OverwriteExisting controls whether existing fields are overwritten.
	OverwriteExisting bool
}

// Injector adds static fields to every log entry that passes through it.
type Injector struct {
	cfg Config
}

// New creates a new Injector with the provided configuration.
func New(cfg Config) *Injector {
	normalised := make(map[string]string, len(cfg.Fields))
	for k, v := range cfg.Fields {
		normalised[strings.ToLower(k)] = v
	}
	cfg.Fields = normalised
	return &Injector{cfg: cfg}
}

// Apply injects the configured static fields into the entry and returns the
// modified copy. The original entry is never mutated.
func (inj *Injector) Apply(entry map[string]any) map[string]any {
	if len(inj.cfg.Fields) == 0 {
		return entry
	}

	out := make(map[string]any, len(entry)+len(inj.cfg.Fields))
	for k, v := range entry {
		out[k] = v
	}

	for field, value := range inj.cfg.Fields {
		existing := parser.HasField(out, field)
		if existing && !inj.cfg.OverwriteExisting {
			continue
		}
		out[field] = value
	}

	return out
}
