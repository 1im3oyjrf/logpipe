package tag

import "strings"

// Rule associates a match condition with a set of tags to inject.
type Rule struct {
	Field string
	Value string
	Tags  []string
}

// Config holds the tagger configuration.
type Config struct {
	Rules           []Rule
	TargetField     string
	CaseInsensitive bool
}

// Tagger injects a tag list into matching log entries.
type Tagger struct {
	cfg Config
}

// New returns a Tagger configured with cfg. If TargetField is empty it
// defaults to "tags".
func New(cfg Config) *Tagger {
	if cfg.TargetField == "" {
		cfg.TargetField = "tags"
	}
	return &Tagger{cfg: cfg}
}

func fieldValue(entry map[string]any, field string, ci bool) string {
	if ci {
		field = strings.ToLower(field)
		for k, v := range entry {
			if strings.ToLower(k) == field {
				if s, ok := v.(string); ok {
					return s
				}
			}
		}
		return ""
	}
	if v, ok := entry[field]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Apply evaluates every rule against entry and injects matching tags into
// the target field. The original entry is never mutated.
func (t *Tagger) Apply(entry map[string]any) map[string]any {
	var matched []string
	for _, r := range t.cfg.Rules {
		v := fieldValue(entry, r.Field, t.cfg.CaseInsensitive)
		ev, cv := r.Value, v
		if t.cfg.CaseInsensitive {
			cv = strings.ToLower(v)
			ev = strings.ToLower(r.Value)
		}
		if cv == ev {
			matched = append(matched, r.Tags...)
		}
	}
	if len(matched) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	out[t.cfg.TargetField] = matched
	return out
}
