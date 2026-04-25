package fold

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls how fields are folded into a single target field.
type Config struct {
	// Fields is the list of source field names to fold.
	Fields []string
	// Target is the destination field name.
	Target string
	// Separator is placed between joined values. Defaults to " ".
	Separator string
	// DropSources removes the source fields after folding.
	DropSources bool
	// CaseInsensitive enables case-insensitive field matching.
	CaseInsensitive bool
}

// Folder collapses multiple log entry fields into a single field.
type Folder struct {
	cfg Config
}

// New creates a Folder from cfg. Returns an error if Target is empty or
// fewer than two source Fields are provided.
func New(cfg Config) (*Folder, error) {
	if strings.TrimSpace(cfg.Target) == "" {
		return nil, fmt.Errorf("fold: target field must not be empty")
	}
	if len(cfg.Fields) < 2 {
		return nil, fmt.Errorf("fold: at least two source fields are required")
	}
	if cfg.Separator == "" {
		cfg.Separator = " "
	}
	return &Folder{cfg: cfg}, nil
}

// Apply folds the configured source fields of entry into the target field and
// returns a new entry. The original entry is never mutated.
func (f *Folder) Apply(entry map[string]any) map[string]any {
	parts := make([]string, 0, len(f.cfg.Fields))
	for _, field := range f.cfg.Fields {
		v := parser.GetString(entry, field, f.cfg.CaseInsensitive)
		if v != "" {
			parts = append(parts, v)
		}
	}
	if len(parts) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	out[f.cfg.Target] = strings.Join(parts, f.cfg.Separator)
	if f.cfg.DropSources {
		for _, field := range f.cfg.Fields {
			if f.cfg.CaseInsensitive {
				for k := range out {
					if strings.EqualFold(k, field) {
						delete(out, k)
					}
				}
			} else {
				delete(out, field)
			}
		}
	}
	return out
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
