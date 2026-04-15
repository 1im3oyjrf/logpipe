package sampler

import (
	"sync/atomic"

	"github.com/your-org/logpipe/internal/reader"
)

// Sampler drops log entries based on a configurable sampling rate.
// A rate of 1 means keep every entry, 10 means keep 1 in every 10, etc.
type Sampler struct {
	rate    uint64
	counter atomic.Uint64
	dropped atomic.Uint64
}

// New returns a Sampler that retains 1 out of every rate entries.
// A rate of 0 or 1 disables sampling (all entries are kept).
func New(rate uint64) *Sampler {
	if rate == 0 {
		rate = 1
	}
	return &Sampler{rate: rate}
}

// Keep returns true if the entry should be forwarded downstream.
// It is safe to call from multiple goroutines.
func (s *Sampler) Keep(_ reader.Entry) bool {
	n := s.counter.Add(1)
	if n%s.rate == 1 {
		return true
	}
	s.dropped.Add(1)
	return false
}

// Dropped returns the total number of entries that have been dropped.
func (s *Sampler) Dropped() uint64 {
	return s.dropped.Load()
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() uint64 {
	return s.rate
}
