// Package window provides a sliding time-window counter for log entries.
// It tracks how many entries pass through within a configurable duration,
// enabling rate-aware decisions in the pipeline.
package window

import (
	"sync"
	"time"
)

// Config holds configuration for a sliding window.
type Config struct {
	// Size is the duration of the sliding window. Defaults to 1 minute.
	Size time.Duration
}

// Window is a thread-safe sliding time-window counter.
type Window struct {
	mu      sync.Mutex
	size    time.Duration
	buckets []bucket
}

type bucket struct {
	at    time.Time
	count int64
}

const defaultSize = time.Minute

// New creates a new Window with the given config.
func New(cfg Config) *Window {
	if cfg.Size <= 0 {
		cfg.Size = defaultSize
	}
	return &Window{size: cfg.Size}
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	w.buckets = append(w.buckets, bucket{at: now, count: n})
	w.evict(now)
}

// Count returns the total number of events within the current window.
func (w *Window) Count() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	var total int64
	for _, b := range w.buckets {
		total += b.count
	}
	return total
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = w.buckets[:0]
}

// evict removes buckets older than the window size. Must be called with mu held.
func (w *Window) evict(now time.Time) {
	cutoff := now.Add(-w.size)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
