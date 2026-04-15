// Package aggregator provides time-window based log entry aggregation,
// grouping entries by a configurable key field and counting occurrences
// within a sliding window duration.
package aggregator

import (
	"sync"
	"time"

	"github.com/user/logpipe/internal/reader"
)

// Bucket holds aggregated counts for a single key within a time window.
type Bucket struct {
	Key    string
	Count  int
	First  time.Time
	Last   time.Time
	Level  string
}

// Aggregator groups log entries by a field value within a rolling window.
type Aggregator struct {
	mu      sync.Mutex
	field   string
	window  time.Duration
	buckets map[string]*Bucket
	now     func() time.Time
}

// New creates an Aggregator that groups entries by field within window.
// If field is empty, "level" is used as the default grouping key.
func New(field string, window time.Duration) *Aggregator {
	if field == "" {
		field = "level"
	}
	if window <= 0 {
		window = 10 * time.Second
	}
	return &Aggregator{
		field:   field,
		window:  window,
		buckets: make(map[string]*Bucket),
		now:     time.Now,
	}
}

// Add records a log entry into the appropriate bucket.
// Entries outside the current window reset their bucket.
func (a *Aggregator) Add(entry reader.Entry) {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := entry.Fields[a.field]
	if key == "" {
		key = "(unknown)"
	}

	now := a.now()
	b, ok := a.buckets[key]
	if !ok || now.Sub(b.First) > a.window {
		a.buckets[key] = &Bucket{
			Key:   key,
			Count: 1,
			First: now,
			Last:  now,
			Level: entry.Level,
		}
		return
	}
	b.Count++
	b.Last = now
}

// Snapshot returns a copy of all current buckets.
func (a *Aggregator) Snapshot() []Bucket {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := make([]Bucket, 0, len(a.buckets))
	for _, b := range a.buckets {
		out = append(out, *b)
	}
	return out
}

// Reset clears all buckets.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = make(map[string]*Bucket)
}
