// Package redact provides field-level redaction for structured log entries.
// It supports exact and pattern-based field name matching to replace sensitive
// values before they are written to any output sink.
package redact

import (
	"regexp"
	"strings"

	"github.com/yourusername/logpipe/internal/reader"
)

// Redactor replaces sensitive field values in log entries.
type Redactor struct {
	exact    map[string]struct{}
	patterns []*regexp.Regexp
	mask     string
}

// Config holds the configuration for a Redactor.
type Config struct {
	// Fields lists exact field names (case-insensitive) to redact.
	Fields []string
	// Patterns lists regular expressions matched against field names.
	Patterns []string
	// Mask is the replacement string; defaults to "[REDACTED]".
	Mask string
}

// New creates a Redactor from the given Config.
// An error is returned if any pattern fails to compile.
func New(cfg Config) (*Redactor, error) {
	mask := cfg.Mask
	if mask == "" {
		mask = "[REDACTED]"
	}

	exact := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		exact[strings.ToLower(f)] = struct{}{}
	}

	patterns := make([]*regexp.Regexp, 0, len(cfg.Patterns))
	for _, p := range cfg.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, re)
	}

	return &Redactor{exact: exact, patterns: patterns, mask: mask}, nil
}

// Apply returns a copy of the entry with matching fields replaced by the mask.
func (r *Redactor) Apply(e reader.Entry) reader.Entry {
	out := reader.Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Source:    e.Source,
		Fields:    make(map[string]any, len(e.Fields)),
	}
	for k, v := range e.Fields {
		if r.shouldRedact(k) {
			out.Fields[k] = r.mask
		} else {
			out.Fields[k] = v
		}
	}
	return out
}

func (r *Redactor) shouldRedact(field string) bool {
	if _, ok := r.exact[strings.ToLower(field)]; ok {
		return true
	}
	for _, re := range r.patterns {
		if re.MatchString(field) {
			return true
		}
	}
	return false
}
