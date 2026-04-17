package normalize

import (
	"strings"

	"github.com/your-org/logpipe/internal/parser"
)

// Config controls field and level normalization behaviour.
type Config struct {
	// FieldMap renames incoming field keys to canonical names.
	// Keys are matched case-insensitively.
	FieldMap map[string]string

	// LevelField is the field that holds the log level (default: "level").
	LevelField string
}

// Normalizer rewrites field names and normalises level values.
type Normalizer struct {
	cfg Config
}

// New returns a Normalizer configured with cfg.
func New(cfg Config) *Normalizer {
	if cfg.LevelField == "" {
		cfg.LevelField = "level"
	}
	// Lower-case all keys in the field map once.
	norm := make(map[string]string, len(cfg.FieldMap))
	for k, v := range cfg.FieldMap {
		norm[strings.ToLower(k)] = v
	}
	cfg.FieldMap = norm
	return &Normalizer{cfg: cfg}
}

// Apply returns a new entry with field names and level value normalised.
// The original entry is never mutated.
func (n *Normalizer) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		canonical := k
		if mapped, ok := n.cfg.FieldMap[strings.ToLower(k)]; ok {
			canonical = mapped
		}
		out[canonical] = v
	}

	// Normalise level to lower-case string.
	if raw, ok := out[n.cfg.LevelField]; ok {
		if s, ok := raw.(string); ok {
			out[n.cfg.LevelField] = strings.ToLower(strings.TrimSpace(s))
		}
	}
	return out
}

// ApplyLevel returns only the normalised level value for the entry,
// or an empty string when the level field is absent.
func (n *Normalizer) ApplyLevel(entry map[string]any) string {
	return strings.ToLower(parser.GetString(entry, n.cfg.LevelField))
}
