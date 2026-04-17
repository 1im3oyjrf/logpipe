package window

import (
	"sync"
	"time"
)

// GuardConfig configures a rate Guard.
type GuardConfig struct {
	// WindowSize is the sliding window duration.
	WindowSize time.Duration
	// Limit is the maximum number of events allowed within the window.
	Limit int64
}

// Guard combines a Window with a threshold limit, allowing callers to check
// whether the current rate exceeds a configured ceiling.
type Guard struct {
	mu     sync.Mutex
	win    *Window
	limit  int64
}

// NewGuard creates a Guard from the given config.
func NewGuard(cfg GuardConfig) *Guard {
	if cfg.Limit <= 0 {
		cfg.Limit = 100
	}
	return &Guard{
		win:   New(Config{Size: cfg.WindowSize}),
		limit: cfg.Limit,
	}
}

// Allow records one event and reports whether the running total is within
// the configured limit. Returns false when the limit has been reached.
func (g *Guard) Allow() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.win.Add(1)
	return g.win.Count() <= g.limit
}

// Count returns the current window count without recording a new event.
func (g *Guard) Count() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.win.Count()
}

// Reset clears the underlying window.
func (g *Guard) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.win.Reset()
}
