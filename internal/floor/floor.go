// Package floor provides a processor that clamps numeric field values to a
// configured minimum, replacing any value below the floor with the floor value.
package floor

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Rule defines a single floor constraint for a named field.
type Rule struct {
	// Field is the log entry field to inspect (case-insensitive).
	Field string
	// Min is the minimum allowed value; values below this are replaced.
	Min float64
}

// Config holds the configuration for the Floor processor.
type Config struct {
	Rules []Rule
}

// Floor replaces numeric field values that fall below a configured minimum.
type Floor struct {
	rules []Rule
}

// New constructs a Floor processor from cfg.
// It returns an error if any rule has an empty field name.
func New(cfg Config) (*Floor, error) {
	rules := make([]Rule, len(cfg.Rules))
	for i, r := range cfg.Rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("floor: rule %d has an empty field name", i)
		}
		rules[i] = Rule{Field: strings.ToLower(r.Field), Min: r.Min}
	}
	return &Floor{rules: rules}, nil
}

// Apply returns a shallow copy of entry with floor constraints applied.
// Fields not present or non-numeric are left unchanged.
func (f *Floor) Apply(entry map[string]any) map[string]any {
	if len(f.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range f.rules {
		for k, v := range out {
			if strings.ToLower(k) != r.Field {
				continue
			}
			val, ok := parser.GetFloat(out, k)
			if !ok {
				continue
			}
			if val < r.Min {
				out[k] = r.Min
			}
			_ = v
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
