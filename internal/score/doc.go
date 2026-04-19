// Package score provides a log-entry scorer that assigns a numeric
// priority value to each entry based on a configurable field and a
// weight map.  The computed score is injected into a configurable
// output field (default: "score") without mutating the original entry.
//
// Example usage:
//
//	s := score.New(score.Config{
//		Field:   "level",
//		Weights: map[string]float64{"error": 10, "warn": 5, "info": 1},
//		Default: 0,
//	})
//	enriched := s.Apply(entry)
package score
