// Package quantize rounds numeric field values to a configured step size,
// reducing cardinality in high-resolution metrics fields.
package quantize

import (
	"fmt"
	"math"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config holds the configuration for the Quantizer.
type Config struct {
	// Field is the entry field to quantize. Required.
	Field string
	// Step is the bucket size. Values are rounded down to the nearest multiple.
	// Must be greater than zero. Defaults to 1.
	Step float64
	// Target is the output field name. Defaults to Field.
	Target string
	// Overwrite controls whether an existing target field is replaced.
	Overwrite bool
}

// Quantizer rounds a numeric field value to the nearest Step multiple.
type Quantizer struct {
	field     string
	step      float64
	target    string
	overwrite bool
}

// New creates a Quantizer from cfg. Returns an error if Field is empty or
// Step is negative.
func New(cfg Config) (*Quantizer, error) {
	if cfg.Field == "" {
		return nil, fmt.Errorf("quantize: field must not be empty")
	}
	if cfg.Step < 0 {
		return nil, fmt.Errorf("quantize: step must be >= 0")
	}
	if cfg.Step == 0 {
		cfg.Step = 1
	}
	target := cfg.Target
	if target == "" {
		target = cfg.Field
	}
	return &Quantizer{
		field:     cfg.Field,
		step:      cfg.Step,
		target:    target,
		overwrite: cfg.Overwrite,
	}, nil
}

// Apply returns a copy of entry with the numeric field quantized. If the field
// is missing or non-numeric the entry is returned unchanged.
func (q *Quantizer) Apply(entry map[string]any) map[string]any {
	v, ok := parser.GetFloat(entry, q.field)
	if !ok {
		return entry
	}
	quantized := math.Floor(v/q.step) * q.step

	out := shallowCopy(entry)
	if _, exists := out[q.target]; exists && !q.overwrite {
		return out
	}
	out[q.target] = quantized
	return out
}

func shallowCopy(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
