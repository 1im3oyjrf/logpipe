// Package dedup provides a deduplication filter for log entries based on
// a configurable window of recently seen message hashes.
package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

// Filter deduplicates log entries by hashing their message and level fields.
// Entries seen within the TTL window are considered duplicates and dropped.
type Filter struct {
	mu    sync.Mutex
	seen  map[string]time.Time
	ttl   time.Duration
	now   func() time.Time
}

// New creates a new deduplication Filter with the given TTL.
// Entries with identical (level, message) pairs within the TTL are suppressed.
func New(ttl time.Duration) *Filter {
	return &Filter{
		seen: make(map[string]time.Time),
		ttl:  ttl,
		now:  time.Now,
	}
}

// IsDuplicate returns true if an equivalent entry was seen within the TTL window.
// It also records the entry if it is not a duplicate.
func (f *Filter) IsDuplicate(entry reader.Entry) bool {
	key := hashEntry(entry)

	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.now()
	f.evict(now)

	if _, exists := f.seen[key]; exists {
		return true
	}

	f.seen[key] = now
	return false
}

// evict removes expired entries from the seen map. Must be called with f.mu held.
func (f *Filter) evict(now time.Time) {
	for key, ts := range f.seen {
		if now.Sub(ts) >= f.ttl {
			delete(f.seen, key)
		}
	}
}

// Size returns the number of entries currently tracked in the dedup window.
func (f *Filter) Size() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}

// Reset clears all tracked entries from the dedup window, effectively
// resetting the filter state. Subsequent entries will not be considered
// duplicates of anything seen before the reset.
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[string]time.Time)
}

func hashEntry(entry reader.Entry) string {
	h := sha256.New()
	h.Write([]byte(entry.Level))
	h.Write([]byte("\x00"))
	h.Write([]byte(entry.Message))
	return hex.EncodeToString(h.Sum(nil))
}
