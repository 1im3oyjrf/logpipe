// Package bound provides a processor that enforces minimum and maximum
// numeric bounds on a named field, clamping values that fall outside the
// configured range and optionally injecting a boolean flag field to
// indicate that clamping occurred.
package bound

import (
	"fmt"
	"strconv"

	"github.com/your-org/logpipe/internal/parser"
)

// Config holds the configuration for the Bound processor.
type Config struct {
	// Field is the entry field to evaluate (required).
	Field string
	// Min is the inclusive lower bound. Nil means no lower bound.
	Min *float64
	// Max is the inclusive upper bound. Nil means no upper bound.
	Max *float64
	// FlagField, when non-empty, receives "true" if the value was clamped.
	FlagField string
}

// Bound clamps a numeric field to a configured range.
type Bound struct {
	cfg Config
}

// New returns a new Bound processor. Returns an error when Field is empty
// or neither Min nor Max is provided.
func New(cfg Config) (*Bound, error) {
	if cfg.Field == "" {
		return nil, fmt.Errorf("bound: field must not be empty")
	}
	if cfg.Min == nil && cfg.Max == nil {
		return nil, fmt.Errorf("bound: at least one of Min or Max must be set")
	}
	if cfg.Min != nil && cfg.Max != nil && *cfg.Min > *cfg.Max {
		return nil, fmt.Errorf("bound: min (%v) must not exceed max (%v)", *cfg.Min, *cfg.Max)
	}
	return &Bound{cfg: cfg}, nil
}

// Apply returns a shallow copy of entry with the target field clamped to the
// configured range. If the field is absent or non-numeric the entry is
// returned unchanged. If FlagField is set it is injected with "true" when
// clamping occurs, or "false" otherwise.
func (b *Bound) Apply(entry map[string]any) map[string]any {
	raw, ok := parser.GetFloat(entry, b.cfg.Field)
	if !ok {
		return entry
	}

	clamped := raw
	if b.cfg.Min != nil && clamped < *b.cfg.Min {
		clamped = *b.cfg.Min
	}
	if b.cfg.Max != nil && clamped > *b.cfg.Max {
		clamped = *b.cfg.Max
	}

	out := shallowCopy(entry)
	out[b.cfg.Field] = clamped

	if b.cfg.FlagField != "" {
		out[b.cfg.FlagField] = strconv.FormatBool(clamped != raw)
	}
	return out
}

func shallowCopy(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
