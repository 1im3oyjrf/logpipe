// Package score assigns a numeric priority score to log entries
// based on configurable field weights. Higher scores indicate
// higher-priority entries for downstream filtering or sorting.
package score

import (
	"strings"

	"github.com/logpipe/logpipe/internal/parser"
)

// Config holds the scoring configuration.
type Config struct {
	// Field is the entry field whose value is matched against Weights.
	Field string
	// Weights maps field values (case-insensitive) to numeric scores.
	Weights map[string]float64
	// Default is the score assigned when no weight matches.
	Default float64
	// OutputField is the field written with the computed score.
	OutputField string
}

// Scorer computes and injects a priority score into each log entry.
type Scorer struct {
	cfg Config
}

// New creates a Scorer from cfg. If OutputField is empty it defaults to "score".
func New(cfg Config) *Scorer {
	if cfg.OutputField == "" {
		cfg.OutputField = "score"
	}
	weights := make(map[string]float64, len(cfg.Weights))
	for k, v := range cfg.Weights {
		weights[strings.ToLower(k)] = v
	}
	cfg.Weights = weights
	return &Scorer{cfg: cfg}
}

// Apply returns a shallow copy of entry with the score field injected.
func (s *Scorer) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}

	val := parser.GetString(entry, s.cfg.Field)
	key := strings.ToLower(val)
	score, ok := s.cfg.Weights[key]
	if !ok {
		score = s.cfg.Default
	}
	out[s.cfg.OutputField] = score
	return out
}
