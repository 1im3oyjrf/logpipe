// Package fieldmap provides a lookup structure for remapping log entry
// field names to canonical output names before formatting or forwarding.
package fieldmap

import "strings"

// Config holds the mapping configuration.
type Config struct {
	// Rules maps source field names (case-insensitive) to target field names.
	Rules map[string]string
	// DropUnmapped drops any field that has no mapping rule.
	DropUnmapped bool
}

// Mapper rewrites field names in a log entry's Fields map.
type Mapper struct {
	rules        map[string]string // lower-cased source -> target
	dropUnmapped bool
}

// New creates a Mapper from cfg. Source keys in Rules are lowercased so
// lookups are case-insensitive at apply time.
func New(cfg Config) *Mapper {
	rules := make(map[string]string, len(cfg.Rules))
	for src, dst := range cfg.Rules {
		rules[strings.ToLower(src)] = dst
	}
	return &Mapper{rules: rules, dropUnmapped: cfg.DropUnmapped}
}

// Apply returns a new Fields map with keys rewritten according to the rules.
// The original map is never mutated.
func (m *Mapper) Apply(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		target, ok := m.rules[strings.ToLower(k)]
		switch {
		case ok:
			out[target] = v
		case m.dropUnmapped:
			// discard
		default:
			out[k] = v
		}
	}
	return out
}

// Rules returns a copy of the internal rules map (keyed by lowercased source).
func (m *Mapper) Rules() map[string]string {
	cp := make(map[string]string, len(m.rules))
	for k, v := range m.rules {
		cp[k] = v
	}
	return cp
}
