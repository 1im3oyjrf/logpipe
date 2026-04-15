// Package buffer provides a fixed-size ring buffer for retaining recent log entries.
package buffer

import (
	"sync"

	"github.com/yourorg/logpipe/internal/reader"
)

// Ring is a thread-safe fixed-capacity circular buffer that retains the most
// recent N log entries. Older entries are overwritten once the buffer is full.
type Ring struct {
	mu       sync.RWMutex
	entries  []reader.Entry
	cap      int
	head     int // next write position
	size     int // current number of valid entries
}

// New creates a new Ring buffer with the given capacity.
// Panics if capacity is less than 1.
func New(capacity int) *Ring {
	if capacity < 1 {
		panic("buffer: capacity must be at least 1")
	}
	return &Ring{
		entries: make([]reader.Entry, capacity),
		cap:     capacity,
	}
}

// Push adds an entry to the ring buffer, overwriting the oldest entry if full.
func (r *Ring) Push(e reader.Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries[r.head] = e
	r.head = (r.head + 1) % r.cap
	if r.size < r.cap {
		r.size++
	}
}

// Snapshot returns a copy of all buffered entries in chronological order
// (oldest first).
func (r *Ring) Snapshot() []reader.Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.size == 0 {
		return nil
	}

	out := make([]reader.Entry, r.size)
	start := (r.head - r.size + r.cap) % r.cap
	for i := 0; i < r.size; i++ {
		out[i] = r.entries[(start+i)%r.cap]
	}
	return out
}

// Len returns the current number of entries held in the buffer.
func (r *Ring) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}

// Cap returns the maximum capacity of the buffer.
func (r *Ring) Cap() int {
	return r.cap
}

// Reset clears all entries from the buffer.
func (r *Ring) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.size = 0
}
