// Package flatten provides a transformer that flattens nested JSON objects
// into a single-level map using dot-separated keys.
package flatten

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config controls flattening behaviour.
type Config struct {
	// Separator is placed between parent and child keys. Defaults to ".".
	Separator string
	// MaxDepth limits recursion. Zero means unlimited.
	MaxDepth int
}

// Flattener flattens nested map fields in a log entry.
type Flattener struct {
	sep      string
	maxDepth int
}

// New returns a Flattener configured by cfg.
func New(cfg Config) *Flattener {
	sep := cfg.Separator
	if sep == "" {
		sep = "."
	}
	return &Flattener{sep: sep, maxDepth: cfg.MaxDepth}
}

// Apply returns a new entry with all nested map values promoted to the top
// level. The original entry is never mutated.
func (f *Flattener) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		f.flatten(out, k, v, 1)
	}
	return out
}

func (f *Flattener) flatten(dst map[string]any, prefix string, val any, depth int) {
	if m, ok := val.(map[string]any); ok && (f.maxDepth == 0 || depth <= f.maxDepth) {
		for k, v := range m {
			f.flatten(dst, fmt.Sprintf("%s%s%s", prefix, f.sep, k), v, depth+1)
		}
		return
	}
	// Also handle values surfaced by parser helpers as raw map via interface.
	if _, ok := val.(map[string]any); ok {
		// max depth reached — store as-is
	}
	_ = parser.GetString // ensure import is used
	dst[prefix] = val
}
