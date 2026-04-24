// Package extract provides a processor that promotes nested fields
// to the top level of a log entry using dot-separated key paths.
package extract

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config holds the configuration for the Extractor.
type Config struct {
	// Paths is a list of dot-separated field paths to extract.
	// e.g. "metadata.request.id" extracts entry["metadata"]["request"]["id"]
	// and places it at entry["metadata.request.id"] (or a custom target).
	Paths []string

	// DropSource removes the original nested key after extraction.
	DropSource bool

	// CaseInsensitive controls whether field matching ignores case.
	CaseInsensitive bool
}

// Extractor promotes nested fields to the top level.
type Extractor struct {
	cfg Config
}

// New creates a new Extractor. Returns an error if no paths are configured.
func New(cfg Config) (*Extractor, error) {
	if len(cfg.Paths) == 0 {
		return nil, fmt.Errorf("extract: at least one path is required")
	}
	return &Extractor{cfg: cfg}, nil
}

// Apply extracts the configured nested paths from entry and returns a new
// map with the extracted values promoted to the top level.
func (e *Extractor) Apply(entry map[string]any) map[string]any {
	out := shallowCopy(entry)
	for _, path := range e.cfg.Paths {
		parts := strings.SplitN(path, ".", 2)
		if len(parts) < 2 {
			continue
		}
		topKey := resolveKey(out, parts[0], e.cfg.CaseInsensitive)
		if topKey == "" {
			continue
		}
		nested, ok := out[topKey].(map[string]any)
		if !ok {
			continue
		}
		leafKey := resolveKey(nested, parts[1], e.cfg.CaseInsensitive)
		if leafKey == "" {
			continue
		}
		out[path] = nested[leafKey]
		if e.cfg.DropSource {
			inner := shallowCopy(nested)
			delete(inner, leafKey)
			out[topKey] = inner
		}
	}
	return out
}

func resolveKey(m map[string]any, key string, caseInsensitive bool) string {
	if _, ok := m[key]; ok {
		return key
	}
	if !caseInsensitive {
		return ""
	}
	lower := strings.ToLower(key)
	for k := range m {
		if strings.ToLower(k) == lower {
			return k
		}
	}
	return ""
}

func shallowCopy(m map[string]any) map[string]any {
	_ = parser.Keys // ensure import is used
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
