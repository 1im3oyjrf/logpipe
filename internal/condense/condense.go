// Package condense merges consecutive log entries that share the same
// message and level into a single entry annotated with a repeat count.
package condense

import (
	"fmt"
	"strings"

	"github.com/yourorg/logpipe/internal/parser"
)

// Config controls condensing behaviour.
type Config struct {
	// MaxRepeat is the maximum number of repeats tracked before flushing.
	// Zero means no limit.
	MaxRepeat int
	// CountField is the field name written with the repeat count.
	// Defaults to "_repeat".
	CountField string
}

// Condenser collapses repeated entries.
type Condenser struct {
	cfg        Config
	last       map[string]interface{}
	count      int
}

// New returns a Condenser with the given config.
func New(cfg Config) *Condenser {
	if cfg.CountField == "" {
		cfg.CountField = "_repeat"
	}
	return &Condenser{cfg: cfg}
}

func key(entry map[string]interface{}) string {
	msg := strings.ToLower(fmt.Sprintf("%v", parser.GetString(entry, "message")))
	lvl := strings.ToLower(fmt.Sprintf("%v", parser.GetString(entry, "level")))
	return lvl + "|" + msg
}

// Apply processes entry. It returns (nil, false) while accumulating repeats
// and (entry, true) when a new distinct entry is seen (flushing the previous).
func (c *Condenser) Apply(entry map[string]interface{}) (map[string]interface{}, bool) {
	if c.last == nil {
		c.last = entry
		c.count = 1
		return nil, false
	}

	if key(entry) == key(c.last) {
		c.count++
		if c.cfg.MaxRepeat > 0 && c.count >= c.cfg.MaxRepeat {
			return c.flush(entry), true
		}
		return nil, false
	}

	out := c.flush(entry)
	return out, true
}

// Flush returns the buffered entry (if any) and resets state.
func (c *Condenser) Flush() (map[string]interface{}, bool) {
	if c.last == nil {
		return nil, false
	}
	out := c.annotate(c.last)
	c.last = nil
	c.count = 0
	return out, true
}

func (c *Condenser) flush(next map[string]interface{}) map[string]interface{} {
	out := c.annotate(c.last)
	c.last = next
	c.count = 1
	return out
}

func (c *Condenser) annotate(entry map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{}, len(entry)+1)
	for k, v := range entry {
		copy[k] = v
	}
	if c.count > 1 {
		copy[c.cfg.CountField] = c.count
	}
	return copy
}
