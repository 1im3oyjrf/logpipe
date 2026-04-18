// Package typecast coerces log entry field values to target types.
package typecast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Rule describes a single field-to-type coercion.
type Rule struct {
	Field  string
	Target string // "string", "int", "float", "bool"
}

// Config holds the set of coercion rules.
type Config struct {
	Rules []Rule
}

// Caster applies type coercions to log entries.
type Caster struct {
	rules []Rule
}

// New returns a Caster configured with the given rules.
func New(cfg Config) *Caster {
	rules := make([]Rule, len(cfg.Rules))
	for i, r := range cfg.Rules {
		rules[i] = Rule{Field: strings.ToLower(r.Field), Target: strings.ToLower(r.Target)}
	}
	return &Caster{rules: rules}
}

// Apply returns a copy of entry with coerced fields.
func (c *Caster) Apply(entry map[string]any) map[string]any {
	if len(c.rules) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range c.rules {
		raw, ok := parser.GetString(entry, r.Field)
		if !ok {
			continue
		}
		coerced, err := coerce(raw, r.Target)
		if err != nil {
			continue
		}
		// find the actual key casing in the map
		for k := range out {
			if strings.ToLower(k) == r.Field {
				out[k] = coerced
				break
			}
		}
	}
	return out
}

func coerce(raw, target string) (any, error) {
	switch target {
	case "string":
		return raw, nil
	case "int":
		return strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	case "float":
		return strconv.ParseFloat(strings.TrimSpace(raw), 64)
	case "bool":
		return strconv.ParseBool(strings.TrimSpace(raw))
	default:
		return nil, fmt.Errorf("unknown target type %q", target)
	}
}
