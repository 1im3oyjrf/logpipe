// Package ceiling provides a processor that enforces a maximum number of
// log entries emitted per time window, dropping entries that exceed the cap.
package ceiling

import (
	"sync"
	"time"

	"github.com/your-org/logpipe/internal/parser"
)

// Config controls the ceiling processor.
type Config struct {
	// Max is the maximum number of entries allowed per Window.
	Max int
	// Window is the rolling time window. Defaults to 1 second.
	Window time.Duration
}

// Ceiling drops entries once Max has been reached within the current Window.
type Ceiling struct {
	mu      sync.Mutex
	cfg     Config
	count   int
	windowEnd time.Time
	dropped int64
}

const defaultWindow = time.Second

// New returns a Ceiling processor. Panics if Max < 1.
func New(cfg Config) *Ceiling {
	if cfg.Max < 1 {
		panic("ceiling: Max must be >= 1")
	}
	if cfg.Window <= 0 {
		cfg.Window = defaultWindow
	}
	return &Ceiling{
		cfg:       cfg,
		windowEnd: time.Now().Add(cfg.Window),
	}
}

// Allow returns true if the entry should be forwarded, false if it should be
// dropped because the cap for the current window has been reached.
func (c *Ceiling) Allow(_ map[string]interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if now.After(c.windowEnd) {
		c.count = 0
		c.windowEnd = now.Add(c.cfg.Window)
	}
	if c.count >= c.cfg.Max {
		c.dropped++
		return false
	}
	c.count++
	return true
}

// Dropped returns the total number of entries dropped since creation.
func (c *Ceiling) Dropped() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.dropped
}

// Reset clears the current window counter and dropped total.
func (c *Ceiling) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = 0
	c.dropped = 0
	c.windowEnd = time.Now().Add(c.cfg.Window)
}

// Apply processes a stream of entries, forwarding only those that pass Allow.
func (c *Ceiling) Apply(in <-chan map[string]interface{}) <-chan map[string]interface{} {
	out := make(chan map[string]interface{})
	go func() {
		defer close(out)
		for entry := range in {
			if c.Allow(parser.ShallowCopy(entry)) {
				out <- entry
			}
		}
	}()
	return out
}
