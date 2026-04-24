// Package abbrev truncates string fields to a maximum length and appends a
// configurable suffix (default "...") so that long values remain readable in
// columnar output without blowing the line budget.
package abbrev

import (
	"strings"

	"github.com/yourorg/logpipe/internal/reader"
)

// Config controls which fields are abbreviated and how.
type Config struct {
	// Fields lists the field names to abbreviate. When empty every string
	// field in the entry is a candidate.
	Fields []string
	// MaxLen is the maximum rune length of a value before it is truncated.
	// Zero or negative values default to 80.
	MaxLen int
	// Suffix is appended to truncated values. Defaults to "...".
	Suffix string
	// CaseInsensitive controls whether field-name matching ignores case.
	CaseInsensitive bool
}

// Abbreviator applies abbreviation rules to log entries.
type Abbreviator struct {
	cfg    Config
	fields map[string]struct{}
}

const defaultMaxLen = 80
const defaultSuffix = "..."

// New creates an Abbreviator from cfg.
func New(cfg Config) *Abbreviator {
	if cfg.MaxLen <= 0 {
		cfg.MaxLen = defaultMaxLen
	}
	if cfg.Suffix == "" {
		cfg.Suffix = defaultSuffix
	}

	set := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		key := f
		if cfg.CaseInsensitive {
			key = strings.ToLower(f)
		}
		set[key] = struct{}{}
	}

	return &Abbreviator{cfg: cfg, fields: set}
}

// Apply returns a shallow copy of entry with matching string fields abbreviated.
func (a *Abbreviator) Apply(entry reader.Entry) reader.Entry {
	out := make(reader.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	for k, v := range out {
		if !a.targeted(k) {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		runes := []rune(s)
		if len(runes) > a.cfg.MaxLen {
			out[k] = string(runes[:a.cfg.MaxLen]) + a.cfg.Suffix
		}
	}
	return out
}

func (a *Abbreviator) targeted(field string) bool {
	if len(a.fields) == 0 {
		return true
	}
	key := field
	if a.cfg.CaseInsensitive {
		key = strings.ToLower(field)
	}
	_, ok := a.fields[key]
	return ok
}
