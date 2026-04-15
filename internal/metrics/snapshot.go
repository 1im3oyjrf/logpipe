package metrics

import "time"

// Snapshot holds an immutable point-in-time view of pipeline metrics.
type Snapshot struct {
	LinesRead     int64
	LinesMatched  int64
	LinesDropped  int64
	Errors        int64
	Uptime        time.Duration
	CapturedAt    time.Time
}

// MatchRate returns the fraction of read lines that matched the filter.
// Returns 0 if no lines have been read.
func (s Snapshot) MatchRate() float64 {
	if s.LinesRead == 0 {
		return 0
	}
	return float64(s.LinesMatched) / float64(s.LinesRead)
}

// DropRate returns the fraction of read lines that were dropped.
// Returns 0 if no lines have been read.
func (s Snapshot) DropRate() float64 {
	if s.LinesRead == 0 {
		return 0
	}
	return float64(s.LinesDropped) / float64(s.LinesRead)
}
