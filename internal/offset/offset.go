package offset

import (
	"encoding/json"
	"strings"
)

// Config controls the offset processor.
type Config struct {
	// Field is the source field whose numeric value is shifted.
	Field string
	// By is the amount added to the field value (may be negative).
	By float64
	// Target is the destination field. Defaults to Field when empty.
	Target string
	// CaseInsensitive enables case-insensitive field matching.
	CaseInsensitive bool
}

// Processor shifts a numeric field by a fixed amount.
type Processor struct {
	cfg Config
}

// New returns a Processor configured by cfg.
func New(cfg Config) *Processor {
	if cfg.Target == "" {
		cfg.Target = cfg.Field
	}
	return &Processor{cfg: cfg}
}

// Apply returns a copy of entry with the numeric field shifted by cfg.By.
// If the field is missing or not numeric the entry is returned unchanged.
func (p *Processor) Apply(entry map[string]any) map[string]any {
	if p.cfg.Field == "" {
		return entry
	}

	key := p.resolveKey(entry)
	if key == "" {
		return entry
	}

	v, ok := entry[key]
	if !ok {
		return entry
	}

	f, ok := toFloat(v)
	if !ok {
		return entry
	}

	out := shallowCopy(entry)
	out[p.cfg.Target] = f + p.cfg.By
	if p.cfg.Target != key {
		delete(out, key)
	}
	return out
}

func (p *Processor) resolveKey(entry map[string]any) string {
	if _, ok := entry[p.cfg.Field]; ok {
		return p.cfg.Field
	}
	if !p.cfg.CaseInsensitive {
		return ""
	}
	lower := strings.ToLower(p.cfg.Field)
	for k := range entry {
		if strings.ToLower(k) == lower {
			return k
		}
	}
	return ""
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
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
