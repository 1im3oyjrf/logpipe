// Package classify assigns a category label to log entries based on
// field value patterns. Rules are evaluated in order; the first match wins.
package classify

import (
	"regexp"
	"strings"

	"github.com/your-org/logpipe/internal/parser"
)

// Rule maps a compiled pattern against a source field to a category string.
type Rule struct {
	Field    string
	Pattern  *regexp.Regexp
	Category string
}

// Config controls the behaviour of the Classifier.
type Config struct {
	// Rules is the ordered list of classification rules.
	Rules []RuleConfig
	// TargetField is the field written with the resolved category.
	// Defaults to "category".
	TargetField string
	// Overwrite controls whether an existing TargetField value is replaced.
	Overwrite bool
}

// RuleConfig is the user-facing (pre-compiled) rule definition.
type RuleConfig struct {
	Field    string
	Pattern  string
	Category string
}

// Classifier applies ordered classification rules to log entries.
type Classifier struct {
	rules       []Rule
	targetField string
	overwrite   bool
}

const defaultTargetField = "category"

// New builds a Classifier from cfg. Returns an error if any pattern fails to
// compile.
func New(cfg Config) (*Classifier, error) {
	tf := strings.TrimSpace(cfg.TargetField)
	if tf == "" {
		tf = defaultTargetField
	}

	rules := make([]Rule, 0, len(cfg.Rules))
	for _, rc := range cfg.Rules {
		re, err := regexp.Compile(rc.Pattern)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{
			Field:    strings.ToLower(strings.TrimSpace(rc.Field)),
			Pattern:  re,
			Category: rc.Category,
		})
	}

	return &Classifier{
		rules:       rules,
		targetField: tf,
		overwrite:   cfg.Overwrite,
	}, nil
}

// Apply evaluates the entry against each rule and injects the category into a
// shallow copy. The original entry is never mutated. If no rule matches the
// entry is returned unchanged.
func (c *Classifier) Apply(entry map[string]any) map[string]any {
	if !c.overwrite {
		if _, exists := parser.HasField(entry, c.targetField); exists {
			return entry
		}
	}

	for _, r := range c.rules {
		val, ok := parser.HasField(entry, r.Field)
		if !ok {
			continue
		}
		str := parser.GetString(val)
		if r.Pattern.MatchString(str) {
			out := shallowCopy(entry)
			out[c.targetField] = r.Category
			return out
		}
	}
	return entry
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
