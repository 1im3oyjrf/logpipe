// Package clamp provides a processor that constrains numeric entry fields
// to a configured [Min, Max] range. Non-numeric fields are passed through
// unchanged. Each rule targets a single field by name (case-insensitive).
package clamp

import (
	"strings"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule defines a clamping constraint for a single field.
type Rule struct {
	// Field is the entry key to clamp (case-insensitive).
	Field string
	// Min is the lower bound; values below it are raised to Min.
	Min *float64
	// Max is the upper bound; values above it are lowered to Max.
	Max *float64
}

// Config holds the set of rules applied by the Clamp processor.
type Config struct {
	Rules []Rule
}

// Clamp applies numeric range constraints to log entry fields.
type Clamp struct {
	rules []Rule
}

// New constructs a Clamp processor from the provided Config.
func New(cfg Config) (*Clamp, error) {
	rules := make([]Rule, len(cfg.Rules))
	for i, r := range cfg.Rules {
		rules[i] = Rule{
			Field: strings.ToLower(r.Field),
			Min:   r.Min,
			Max:   r.Max,
		}
	}
	return &Clamp{rules: rules}, nil
}

// Apply returns a new entry with numeric fields clamped according to the
// configured rules. The original entry is never modified.
func (c *Clamp) Apply(e reader.Entry) reader.Entry {
	if len(c.rules) == 0 {
		return e
	}
	out := shallowCopy(e)
	for _, r := range c.rules {
		matchKey(out, r)
	}
	return out
}

// matchKey locates the target field in the entry (case-insensitive) and
// applies the min/max constraints when the value is numeric.
func matchKey(e reader.Entry, r Rule) {
	for k, v := range e {
		if strings.ToLower(k) != r.Field {
			continue
		}
		f, ok := toFloat(v)
		if !ok {
			return
		}
		if r.Min != nil && f < *r.Min {
			f = *r.Min
		}
		if r.Max != nil && f > *r.Max {
			f = *r.Max
		}
		e[k] = f
		return
	}
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func shallowCopy(e reader.Entry) reader.Entry {
	out := make(reader.Entry, len(e))
	for k, v := range e {
		out[k] = v
	}
	return out
}
