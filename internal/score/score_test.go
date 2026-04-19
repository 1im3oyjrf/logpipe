package score_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/score"
)

func base() map[string]any {
	return map[string]any{"level": "error", "msg": "boom"}
}

func TestApply_MatchingWeight_InjectsScore(t *testing.T) {
	s := score.New(score.Config{
		Field:   "level",
		Weights: map[string]float64{"error": 10, "warn": 5, "info": 1},
	})
	out := s.Apply(base())
	if out["score"] != 10.0 {
		t.Fatalf("expected 10, got %v", out["score"])
	}
}

func TestApply_NoMatch_UsesDefault(t *testing.T) {
	s := score.New(score.Config{
		Field:   "level",
		Weights: map[string]float64{"warn": 5},
		Default: 2,
	})
	out := s.Apply(base()) // level=error, not in weights
	if out["score"] != 2.0 {
		t.Fatalf("expected default 2, got %v", out["score"])
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	s := score.New(score.Config{
		Field:   "level",
		Weights: map[string]float64{"ERROR": 10},
	})
	out := s.Apply(base())
	if out["score"] != 10.0 {
		t.Fatalf("expected 10, got %v", out["score"])
	}
}

func TestApply_CustomOutputField(t *testing.T) {
	s := score.New(score.Config{
		Field:       "level",
		Weights:     map[string]float64{"error": 99},
		OutputField: "priority",
	})
	out := s.Apply(base())
	if _, ok := out["score"]; ok {
		t.Fatal("default 'score' field should not be present")
	}
	if out["priority"] != 99.0 {
		t.Fatalf("expected 99, got %v", out["priority"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	s := score.New(score.Config{
		Field:   "level",
		Weights: map[string]float64{"error": 10},
	})
	orig := base()
	_ = s.Apply(orig)
	if _, ok := orig["score"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}

func TestApply_MissingField_UsesDefault(t *testing.T) {
	s := score.New(score.Config{
		Field:   "severity",
		Weights: map[string]float64{"critical": 100},
		Default: 0,
	})
	out := s.Apply(base()) // no 'severity' field
	if out["score"] != 0.0 {
		t.Fatalf("expected 0, got %v", out["score"])
	}
}
