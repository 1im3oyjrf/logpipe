package metrics

import (
	"sync/atomic"
	"time"
)

// Counters tracks pipeline processing statistics.
type Counters struct {
	LinesRead    atomic.Int64
	LinesMatched atomic.Int64
	LinesDropped atomic.Int64
	ParseErrors  atomic.Int64
	startedAt    time.Time
}

// New returns a new Counters instance with the start time set to now.
func New() *Counters {
	return &Counters{
		startedAt: time.Now(),
	}
}

// IncRead increments the lines-read counter by 1.
func (c *Counters) IncRead() { c.LinesRead.Add(1) }

// IncMatched increments the lines-matched counter by 1.
func (c *Counters) IncMatched() { c.LinesMatched.Add(1) }

// IncDropped increments the lines-dropped counter by 1.
func (c *Counters) IncDropped() { c.LinesDropped.Add(1) }

// IncParseError increments the parse-error counter by 1.
func (c *Counters) IncParseError() { c.ParseErrors.Add(1) }

// Snapshot returns a point-in-time copy of the current counter values.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:    c.LinesRead.Load(),
		LinesMatched: c.LinesMatched.Load(),
		LinesDropped: c.LinesDropped.Load(),
		ParseErrors:  c.ParseErrors.Load(),
		Uptime:       time.Since(c.startedAt).Round(time.Millisecond),
	}
}

// Snapshot holds an immutable copy of counter values at a moment in time.
type Snapshot struct {
	LinesRead    int64
	LinesMatched int64
	LinesDropped int64
	ParseErrors  int64
	Uptime       time.Duration
}

// MatchRate returns the fraction of read lines that were matched, in the
// range [0.0, 1.0]. Returns 0 if no lines have been read.
func (s Snapshot) MatchRate() float64 {
	if s.LinesRead == 0 {
		return 0
	}
	return float64(s.LinesMatched) / float64(s.LinesRead)
}
