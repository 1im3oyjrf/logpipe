package pivot

import (
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls how entries are pivoted.
type Config struct {
	// KeyField is the field whose value becomes a new field name.
	KeyField string
	// ValueField is the field whose value is placed under the new field name.
	ValueField string
	// DropSource removes KeyField and ValueField from the output entry.
	DropSource bool
}

// Pivoter promotes a key/value pair from two fields into a top-level field.
type Pivoter struct {
	cfg Config
}

// New returns a Pivoter configured with cfg.
func New(cfg Config) *Pivoter {
	if cfg.KeyField == "" {
		cfg.KeyField = "key"
	}
	if cfg.ValueField == "" {
		cfg.ValueField = "value"
	}
	return &Pivoter{cfg: cfg}
}

// Apply returns a shallow copy of entry with the pivot applied.
// If KeyField is absent or its value is empty the entry is returned unchanged.
func (p *Pivoter) Apply(entry map[string]any) map[string]any {
	keyVal := parser.GetString(entry, p.cfg.KeyField)
	if keyVal == "" {
		return entry
	}
	valVal, ok := entry[p.findKey(entry, p.cfg.ValueField)]
	if !ok {
		return entry
	}

	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	out[keyVal] = valVal
	if p.cfg.DropSource {
		delete(out, p.findKey(out, p.cfg.KeyField))
		delete(out, p.findKey(out, p.cfg.ValueField))
	}
	return out
}

func (p *Pivoter) findKey(entry map[string]any, name string) string {
	if _, ok := entry[name]; ok {
		return name
	}
	lower := strings.ToLower(name)
	for k := range entry {
		if strings.ToLower(k) == lower {
			return k
		}
	}
	return name
}
