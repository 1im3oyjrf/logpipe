// Package prefix prepends a static string to a named log field.
package prefix

import "strings"

// Config controls which field is prefixed and what value is prepended.
type Config struct {
	// Field is the log entry key to modify (case-insensitive).
	Field string
	// Value is the string prepended to the field's existing value.
	Value string
	// Sep is placed between Value and the original content. Defaults to empty.
	Sep string
}

// Prefixer applies a prefix transformation to log entries.
type Prefixer struct {
	field string
	value string
	sep   string
}

// New creates a Prefixer from cfg. If Field or Value is empty the Prefixer
// is a no-op.
func New(cfg Config) *Prefixer {
	return &Prefixer{
		field: strings.ToLower(cfg.Field),
		value: cfg.Value,
		sep:   cfg.Sep,
	}
}

// Apply returns a shallow copy of entry with the configured field prefixed.
// If the field is absent, or the Prefixer is a no-op, the original map is
// returned unchanged.
func (p *Prefixer) Apply(entry map[string]any) map[string]any {
	if p.field == "" || p.value == "" {
		return entry
	}

	// Locate the key case-insensitively.
	var matched string
	for k := range entry {
		if strings.ToLower(k) == p.field {
			matched = k
			break
		}
	}
	if matched == "" {
		return entry
	}

	existing, ok := entry[matched].(string)
	if !ok {
		return entry
	}

	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	out[matched] = p.value + p.sep + existing
	return out
}
