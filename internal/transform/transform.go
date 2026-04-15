package transform

import (
	"strings"

	"github.com/yourusername/logpipe/internal/reader"
)

// Transformer applies field-level transformations to log entries.
type Transformer struct {
	redactFields map[string]bool
	renameFields map[string]string
	addFields    map[string]string
}

// Config holds the transformation configuration.
type Config struct {
	// RedactFields is a list of field names whose values will be replaced with "[REDACTED]".
	RedactFields []string
	// RenameFields maps old field names to new field names.
	RenameFields map[string]string
	// AddFields maps field names to static values to inject into every entry.
	AddFields map[string]string
}

// New creates a new Transformer from the given Config.
func New(cfg Config) *Transformer {
	redact := make(map[string]bool, len(cfg.RedactFields))
	for _, f := range cfg.RedactFields {
		redact[strings.ToLower(f)] = true
	}
	return &Transformer{
		redactFields: redact,
		renameFields: cfg.RenameFields,
		addFields:    cfg.AddFields,
	}
}

// Apply returns a new log entry with transformations applied.
// The original entry is not modified.
func (t *Transformer) Apply(e reader.LogEntry) reader.LogEntry {
	out := reader.LogEntry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Source:    e.Source,
		Fields:    make(map[string]any, len(e.Fields)+len(t.addFields)),
	}

	// Copy and transform existing fields.
	for k, v := range e.Fields {
		newKey := k
		if renamed, ok := t.renameFields[k]; ok {
			newKey = renamed
		}
		if t.redactFields[strings.ToLower(k)] {
			out.Fields[newKey] = "[REDACTED]"
		} else {
			out.Fields[newKey] = v
		}
	}

	// Inject static fields (do not overwrite existing keys).
	for k, v := range t.addFields {
		if _, exists := out.Fields[k]; !exists {
			out.Fields[k] = v
		}
	}

	return out
}
