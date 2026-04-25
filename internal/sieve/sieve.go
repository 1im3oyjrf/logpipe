// Package sieve provides a probabilistic entry filter using a fixed-size
// bit array (Bloom-filter-inspired) to drop entries whose field value has
// already been seen within the current generation. Unlike [distinct], sieve
// trades perfect accuracy for O(1) memory regardless of cardinality.
package sieve

import (
	"hash/fnv"
	"strings"
	"sync"

	"github.com/logpipe/logpipe/internal/parser"
)

const defaultField = "message"
const defaultSlots = 1024

// Config controls sieve behaviour.
type Config struct {
	// Field is the entry field whose value is hashed. Defaults to "message".
	Field string
	// Slots is the number of bit-slots in the internal array. Must be > 0.
	// Larger values reduce false-positive collision rates.
	Slots int
	// CaseInsensitive lowercases the value before hashing.
	CaseInsensitive bool
}

// Sieve drops entries whose hashed field value collides with a previously
// seen slot. Collisions cause false positives (entries wrongly dropped) but
// never false negatives.
type Sieve struct {
	field           string
	slots           int
	caseInsensitive bool
	mu              sync.Mutex
	bits            []bool
}

// New constructs a Sieve from cfg.
func New(cfg Config) *Sieve {
	if cfg.Field == "" {
		cfg.Field = defaultField
	}
	if cfg.Slots <= 0 {
		cfg.Slots = defaultSlots
	}
	return &Sieve{
		field:           cfg.Field,
		slots:           cfg.Slots,
		caseInsensitive: cfg.CaseInsensitive,
		bits:            make([]bool, cfg.Slots),
	}
}

// Allow returns true when the entry should be forwarded, false when the
// hashed slot was already occupied (probable duplicate).
func (s *Sieve) Allow(entry map[string]any) bool {
	v := parser.GetString(entry, s.field)
	if s.caseInsensitive {
		v = strings.ToLower(v)
	}
	idx := s.slot(v)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bits[idx] {
		return false
	}
	s.bits[idx] = true
	return true
}

// Reset clears all slots, starting a new generation.
func (s *Sieve) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.bits {
		s.bits[i] = false
	}
}

func (s *Sieve) slot(v string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(v))
	return int(h.Sum32()) % s.slots
}
