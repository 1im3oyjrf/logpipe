// Package stencil provides a processor that copies a fixed set of fields
// from a template entry into every log entry that passes through it.  Fields
// already present in the entry are left untouched unless Overwrite is set.
package stencil

import (
	"fmt"
	"strings"

	"github.com/your-org/logpipe/internal/parser"
)

// Config controls which fields are stamped onto each entry.
type Config struct {
	// Fields is the map of field-name → value to inject.
	Fields map[string]any
	// Overwrite replaces existing fields when true.
	Overwrite bool
	// CaseInsensitive matches destination keys without regard to case.
	CaseInsensitive bool
}

// Stencil stamps a fixed set of fields onto every entry.
type Stencil struct {
	fields          map[string]any
	overwrite       bool
	caseInsensitive bool
}

// New returns a Stencil configured from cfg.
// It returns an error when Fields is empty.
func New(cfg Config) (*Stencil, error) {
	if len(cfg.Fields) == 0 {
		return nil, fmt.Errorf("stencil: Fields must not be empty")
	}
	norm := make(map[string]any, len(cfg.Fields))
	for k, v := range cfg.Fields {
		key := k
		if cfg.CaseInsensitive {
			key = strings.ToLower(k)
		}
		norm[key] = v
	}
	return &Stencil{
		fields:          norm,
		overwrite:       cfg.Overwrite,
		caseInsensitive: cfg.CaseInsensitive,
	}, nil
}

// Apply returns a shallow copy of entry with the stencil fields injected.
func (s *Stencil) Apply(entry map[string]any) map[string]any {
	out := shallowCopy(entry)
	for k, v := range s.fields {
		target := k
		if s.caseInsensitive {
			// find the actual key in the entry (case-insensitive)
			for ek := range out {
				if strings.EqualFold(ek, k) {
					target = ek
					break
				}
			}
		}
		if _, exists := out[target]; exists && !s.overwrite {
			continue
		}
		out[target] = v
	}
	_ = parser.Keys // satisfy import
	return out
}

func shallowCopy(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
