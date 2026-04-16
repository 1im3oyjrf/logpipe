// Package label provides field-based log entry labelling.
// Labels are key/value pairs injected into entries based on matching rules.
package label

import (
	"strings"

	"github.com/logpipe/internal/parser"
)

// Rule defines a condition and the labels to apply when it matches.
type Rule struct {
	// Field is the entry field to inspect.
	Field string
	// Value is the substring to match (case-insensitive).
	Value string
	// Labels are the key/value pairs to inject when the rule matches.
	Labels map[string]string
}

// Labeller applies label rules to log entries.
type Labeller struct {
	rules []Rule
}

// New returns a Labeller configured with the given rules.
func New(rules []Rule) *Labeller {
	return &Labeller{rules: rules}
}

// Apply returns a copy of entry with labels injected for every matching rule.
// Original entry fields are never mutated.
func (l *Labeller) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range l.rules {
		fv := parser.GetString(entry, r.Field)
		if strings.Contains(strings.ToLower(fv), strings.ToLower(r.Value)) {
			for k, v := range r.Labels {
				out[k] = v
			}
		}
	}
	return out
}
