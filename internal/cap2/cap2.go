// Package cap2 provides a pipeline stage that limits the rate of entries
// forwarded downstream by dropping entries once a per-field value count
// exceeds a configured maximum within a sliding time window.
package cap2

import (
	"strings"
	"sync"
	"time"

	"github.com/user/logpipe/internal/parser"
)

// Entry is the common log entry type used across the pipeline.
type Entry = map[string]any

// Config controls the behaviour of the Capper.
type Config struct {
	// Field is the entry field whose value is counted (e.g. "level", "source").
	// An empty Field counts all entries under a single shared key.
	Field string
	// Max is the maximum number of entries allowed per unique field value
	// within Window. Zero means no limit.
	Max int
	// Window is the duration of the sliding count window. Defaults to 1 minute.
	Window time.Duration
	// CaseInsensitive controls whether field value comparisons ignore case.
	CaseInsensitive bool
}

type bucket struct {
	count  int
	expiry time.Time
}

// Capper tracks per-value entry counts and drops entries that exceed Max
// within the configured Window.
type Capper struct {
	cfg Config
	mu      sync.Mutex
	buckets map[string]*bucket
}

// New constructs a Capper from cfg. If cfg.Window is zero it defaults to
// one minute.
func New(cfg Config) *Capper {
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	return &Capper{
		cfg:     cfg,
		buckets: make(map[string]*bucket),
	}
}

// Allow returns true when the entry should be forwarded downstream and false
// when it has been suppressed because the per-value count exceeds Max.
func (c *Capper) Allow(e Entry) bool {
	if c.cfg.Max <= 0 {
		return true
	}

	key := c.keyFor(e)

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	b, ok := c.buckets[key]
	if !ok || now.After(b.expiry) {
		c.buckets[key] = &bucket{count: 1, expiry: now.Add(c.cfg.Window)}
		return true
	}
	b.count++
	return b.count <= c.cfg.Max
}

// Reset clears all tracked counters.
func (c *Capper) Reset() {
	c.mu.Lock()
	c.buckets = make(map[string]*bucket)
	c.mu.Unlock()
}

func (c *Capper) keyFor(e Entry) string {
	if c.cfg.Field == "" {
		return "__all__"
	}
	v := parser.GetString(e, c.cfg.Field)
	if c.cfg.CaseInsensitive {
		v = strings.ToLower(v)
	}
	return v
}
