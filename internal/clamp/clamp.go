// Package clamp provides a processor that constrains numeric field values
// to a configured [min, max] range, clamping any out-of-range value to the
// nearest boundary.
package clamp

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Rule describes a single field clamping rule.
type Rule struct {
	Field string
	Min   float64
	Max   float64
}

// Config holds the configuration for the Clamp processor.
type Config struct {
	Rules           []Rule
	CaseInsensitive bool
}

// Clamp constrains numeric fields to [Min, Max].
type Clamp struct {
	cfg Config
}

// New returns a Clamp processor or an error if any rule is invalid.
func New(cfg Config) (*Clamp, error) {
	for i, r := range cfg.Rules {
		if r.Field == "" {
			return nil, fmt.Errorf("clamp: rule %d has empty field", i)
		}
		if r.Min > r.Max {
			return nil, fmt.Errorf("clamp: rule %d min (%v) exceeds max (%v)", i, r.Min, r.Max)
		}
	}
	return &Clamp{cfg: cfg}, nil
}

// Apply returns a shallow copy of entry with numeric fields clamped per rules.
func (c *Clamp) Apply(entry map[string]any) map[string]any {
	if len(c.cfg.Rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, rule := range c.cfg.Rules {
		key := matchKey(out, rule.Field, c.cfg.CaseInsensitive)
		if key == "" {
			continue
		}
		v, ok := parser.GetFloat(out, key)
		if !ok {
			continue
		}
		if v < rule.Min {
			v = rule.Min
		} else if v > rule.Max {
			v = rule.Max
		}
		out[key] = v
	}
	return out
}

func matchKey(entry map[string]any, field string, ci bool) string {
	if !ci {
		if _, ok := entry[field]; ok {
			return field
		}
		return ""
	}
	lower := strings.ToLower(field)
	for k := range entry {
		if strings.ToLower(k) == lower {
			return k
		}
	}
	return ""
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
