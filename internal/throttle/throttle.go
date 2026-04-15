// Package throttle provides a per-key throttling mechanism that suppresses
// repeated log entries within a configurable time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last emission time per key and suppresses entries
// that arrive within the cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
	now      func() time.Time
}

// New creates a Throttle with the given cooldown duration.
// Entries sharing the same key will be suppressed if they arrive within
// cooldown of the previous emission.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the entry identified by key should be forwarded.
// It returns false when the key was seen within the cooldown window.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.seen[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.seen[key] = now
	return true
}

// Reset removes all tracked keys, allowing all entries through again.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.seen = make(map[string]time.Time)
}

// Len returns the number of keys currently being tracked.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.seen)
}

// Evict removes entries whose last-seen time is older than the cooldown,
// freeing memory for long-running processes.
func (t *Throttle) Evict() {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	for k, last := range t.seen {
		if now.Sub(last) >= t.cooldown {
			delete(t.seen, k)
		}
	}
}
