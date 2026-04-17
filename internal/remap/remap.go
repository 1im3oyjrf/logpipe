// Package remap provides value-level remapping for log entry fields.
// It replaces specific field values with configured substitutions,
// useful for normalising codes, status strings, or legacy values.
package remap

import (
	"strings"

	"github.com/yourorg/logpipe/internal/parser"
)

// Rule maps a source field's value to a replacement.
type Rule struct {
	Field       string
	From        string
	To          string
	CaseSensitive bool
}

// Config holds configuration for the Remapper.
type Config struct {
	Rules []Rule
}

// Remapper applies value substitutions to log entries.
type Remapper struct {
	rules []Rule
}

// New creates a Remapper from the given Config.
func New(cfg Config) *Remapper {
	rules := make([]Rule, len(cfg.Rules))
	copy(rules, cfg.Rules)
	return &Remapper{rules: rules}
}

// Apply returns a new entry with matching field values replaced.
// The original entry is never mutated.
func (r *Remapper) Apply(entry map[string]any) map[string]any {
	out := shallowCopy(entry)
	for _, rule := range r.rules {
		val, ok := parser.GetString(out, rule.Field)
		if !ok {
			continue
		}
		matches := false
		if rule.CaseSensitive {
			matches = val == rule.From
		} else {
			matches = strings.EqualFold(val, rule.From)
		}
		if matches {
			out[rule.Field] = rule.To
		}
	}
	return out
}

func shallowCopy(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
