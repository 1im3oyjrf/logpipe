// Package split provides a processor that splits a single log entry into
// multiple entries by expanding a repeated/array field into individual records.
package split

import (
	"fmt"
	"strings"

	"github.com/your-org/logpipe/internal/parser"
)

// Entry mirrors the shared log entry type used across logpipe.
type Entry = map[string]any

// Config controls which field to expand.
type Config struct {
	// Field is the entry key whose value is a []any to expand.
	Field string
	// TargetField is the key written on each child entry.
	// Defaults to Field when empty.
	TargetField string
}

// Splitter expands array-valued fields into individual entries.
type Splitter struct {
	field  string
	target string
}

// New returns a Splitter for cfg.
func New(cfg Config) (*Splitter, error) {
	f := strings.TrimSpace(cfg.Field)
	if f == "" {
		return nil, fmt.Errorf("split: field must not be empty")
	}
	t := strings.TrimSpace(cfg.TargetField)
	if t == "" {
		t = f
	}
	return &Splitter{field: f, target: t}, nil
}

// Apply returns one entry per element in the array field.
// If the field is absent or not a slice the original entry is returned as-is
// in a single-element slice.
func (s *Splitter) Apply(e Entry) []Entry {
	key := parser.Keys(e)
	// find case-insensitive match
	actual := ""
	for _, k := range key {
		if strings.EqualFold(k, s.field) {
			actual = k
			break
		}
	}
	if actual == "" {
		return []Entry{e}
	}

	raw, ok := e[actual].([]any)
	if !ok || len(raw) == 0 {
		return []Entry{e}
	}

	out := make([]Entry, 0, len(raw))
	for _, v := range raw {
		child := make(Entry, len(e))
		for k, val := range e {
			if !strings.EqualFold(k, s.field) {
				child[k] = val
			}
		}
		child[s.target] = v
		out = append(out, child)
	}
	return out
}
