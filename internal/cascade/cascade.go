package cascade

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config holds the configuration for the cascade processor.
type Config struct {
	// Rules is an ordered list of field-value conditions paired with the
	// target field and value to inject when the condition matches.
	Rules []Rule
	// StopOnFirst causes evaluation to stop after the first matching rule.
	// Defaults to true.
	StopOnFirst bool
}

// Rule represents a single cascade condition and its resulting action.
type Rule struct {
	Field  string
	Value  string
	Target string
	Set    string
}

// Cascade evaluates an ordered list of field-matching rules and injects a
// value into a target field for the first (or all) matching rules.
type Cascade struct {
	rules       []Rule
	stopOnFirst bool
}

// New constructs a Cascade from cfg. Returns an error if any rule is
// missing a required field.
func New(cfg Config) (*Cascade, error) {
	for i, r := range cfg.Rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("cascade: rule %d: field must not be empty", i)
		}
		if strings.TrimSpace(r.Target) == "" {
			return nil, fmt.Errorf("cascade: rule %d: target must not be empty", i)
		}
	}
	stop := true
	if !cfg.StopOnFirst {
		stop = false
	}
	rules := make([]Rule, len(cfg.Rules))
	copy(rules, cfg.Rules)
	return &Cascade{rules: rules, stopOnFirst: stop}, nil
}

// Apply evaluates each rule against entry and returns a new entry with the
// appropriate target fields injected. The original entry is never mutated.
func (c *Cascade) Apply(entry map[string]any) map[string]any {
	out := shallowCopy(entry)
	for _, r := range c.rules {
		v := parser.GetString(out, r.Field)
		if !strings.EqualFold(v, r.Value) {
			continue
		}
		out[r.Target] = r.Set
		if c.stopOnFirst {
			break
		}
	}
	return out
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
