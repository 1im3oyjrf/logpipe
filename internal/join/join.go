package join

import (
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls how fields are joined into a new field.
type Config struct {
	// Fields is the ordered list of source field names to join.
	Fields []string
	// Separator is placed between each field value.
	Separator string
	// Target is the destination field name for the joined value.
	Target string
	// DropSources removes the source fields after joining.
	DropSources bool
}

// Joiner combines multiple log entry fields into a single new field.
type Joiner struct {
	cfg Config
}

// New returns a Joiner configured with cfg.
// If Target is empty it defaults to "joined".
// If Separator is empty it defaults to " ".
func New(cfg Config) *Joiner {
	if cfg.Target == "" {
		cfg.Target = "joined"
	}
	if cfg.Separator == "" {
		cfg.Separator = " "
	}
	return &Joiner{cfg: cfg}
}

// Apply returns a shallow copy of entry with the joined field injected.
// Source fields that are missing or non-string are skipped.
func (j *Joiner) Apply(entry map[string]any) map[string]any {
	if len(j.cfg.Fields) == 0 {
		return entry
	}

	parts := make([]string, 0, len(j.cfg.Fields))
	for _, f := range j.cfg.Fields {
		v := parser.GetString(entry, f)
		if v != "" {
			parts = append(parts, v)
		}
	}

	out := shallowCopy(entry)
	out[j.cfg.Target] = strings.Join(parts, j.cfg.Separator)

	if j.cfg.DropSources {
		for _, f := range j.cfg.Fields {
			for k := range out {
				if strings.EqualFold(k, f) {
					delete(out, k)
				}
			}
		}
	}

	return out
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
