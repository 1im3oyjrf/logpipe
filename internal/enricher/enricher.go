package enricher

import (
	"strings"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

// Config holds static fields to inject into every log entry.
type Config struct {
	// StaticFields are key/value pairs added to every entry.
	StaticFields map[string]string
	// AddTimestamp injects an "enriched_at" field with the current UTC time.
	AddTimestamp bool
	// HostField, if non-empty, sets that field to the system hostname.
	HostField string
}

// Enricher adds static or derived fields to log entries.
type Enricher struct {
	cfg    Config
	hostname string
}

// New creates an Enricher from cfg. Hostname resolution is best-effort.
func New(cfg Config) *Enricher {
	host := ""
	if cfg.HostField != "" {
		// avoid importing os at call sites; resolve once here
		host = resolveHostname()
	}
	return &Enricher{cfg: cfg, hostname: host}
}

// Apply returns a copy of entry with configured fields injected.
func (e *Enricher) Apply(entry reader.Entry) reader.Entry {
	out := entry
	if out.Fields == nil {
		out.Fields = make(map[string]any)
	} else {
		// shallow-copy so we don't mutate the original map
		copy := make(map[string]any, len(entry.Fields))
		for k, v := range entry.Fields {
			copy[k] = v
		}
		out.Fields = copy
	}

	for k, v := range e.cfg.StaticFields {
		out.Fields[k] = v
	}

	if e.cfg.AddTimestamp {
		out.Fields["enriched_at"] = time.Now().UTC().Format(time.RFC3339)
	}

	if e.cfg.HostField != "" && e.hostname != "" {
		out.Fields[strings.ToLower(e.cfg.HostField)] = e.hostname
	}

	return out
}

func resolveHostname() string {
	import_os_once.Do(func() {})
	return hostnameOnce()
}
