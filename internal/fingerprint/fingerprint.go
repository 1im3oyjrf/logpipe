// Package fingerprint computes a stable string identity for a log entry
// based on a configurable set of fields. Entries that share the same
// fingerprint can be grouped, deduplicated, or correlated downstream.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/logpipe/internal/parser"
)

// Config holds the options for the Fingerprinter.
type Config struct {
	// Fields is the ordered list of field names to include in the fingerprint.
	// When empty all fields are used, sorted alphabetically.
	Fields []string

	// OutputField is the entry field where the fingerprint is written.
	// Defaults to "_fp".
	OutputField string

	// CaseInsensitive controls whether field name lookup ignores case.
	CaseInsensitive bool
}

// Fingerprinter computes and injects a fingerprint into each log entry.
type Fingerprinter struct {
	cfg Config
}

// New creates a Fingerprinter from cfg, applying defaults where needed.
func New(cfg Config) *Fingerprinter {
	if cfg.OutputField == "" {
		cfg.OutputField = "_fp"
	}
	return &Fingerprinter{cfg: cfg}
}

// Apply computes the fingerprint for entry and returns a shallow copy with
// the fingerprint written to the configured output field.
func (f *Fingerprinter) Apply(entry map[string]any) map[string]any {
	keys := f.resolveKeys(entry)

	var sb strings.Builder
	for _, k := range keys {
		v := parser.GetString(entry, k)
		fmt.Fprintf(&sb, "%s=%s;", k, v)
	}

	sum := sha256.Sum256([]byte(sb.String()))
	fp := hex.EncodeToString(sum[:8]) // 16-char prefix is sufficient

	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	out[f.cfg.OutputField] = fp
	return out
}

// resolveKeys returns the sorted list of field names to hash.
func (f *Fingerprinter) resolveKeys(entry map[string]any) []string {
	if len(f.cfg.Fields) > 0 {
		result := make([]string, 0, len(f.cfg.Fields))
		for _, want := range f.cfg.Fields {
			for k := range entry {
				if f.match(k, want) {
					result = append(result, k)
					break
				}
			}
		}
		return result
	}

	keys := make([]string, 0, len(entry))
	for k := range entry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (f *Fingerprinter) match(a, b string) bool {
	if f.cfg.CaseInsensitive {
		return strings.EqualFold(a, b)
	}
	return a == b
}
