// Package coalesce provides a transformer that merges a prioritised list of
// fields into a single canonical field, using the first non-empty value found.
package coalesce

import (
	"strings"

	"github.com/yourorg/logpipe/internal/parser"
)

// Rule describes one coalesce operation.
type Rule struct {
	// Sources is an ordered list of field names to try.
	Sources []string
	// Target is the field that receives the first non-empty value.
	Target string
	// KeepSources, when false, removes the source fields after merging.
	KeepSources bool
}

// Config holds all rules for the transformer.
type Config struct {
	Rules []Rule
}

// Transformer applies coalesce rules to log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer configured with cfg.
func New(cfg Config) *Transformer {
	return &Transformer{rules: cfg.Rules}
}

// Apply returns a new entry with coalesce rules applied.
// The original entry is never mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	out := shallowCopy(entry)
	for _, rule := range t.rules {
		value := ""
		picked := ""
		for _, src := range rule.Sources {
			v := parser.GetString(out, src)
			if strings.TrimSpace(v) != "" {
				value = v
				picked = src
				break
			}
		}
		if value == "" {
			continue
		}
		out[rule.Target] = value
		if !rule.KeepSources {
			for _, src := range rule.Sources {
				if src != rule.Target && src != picked {
					delete(out, src)
				}
				if src == picked && src != rule.Target {
					delete(out, src)
				}
			}
		}
	}
	return out
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
