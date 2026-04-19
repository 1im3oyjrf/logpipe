package clamp_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/clamp"
)

func base() map[string]any {
	return map[string]any{
		"level":    "info",
		"message":  "hello",
		"duration": float64(120),
		"score":    float64(5),
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	c, _ := clamp.New(clamp.Config{})
	out := c.Apply(base())
	if out["duration"] != float64(120) {
		t.Fatalf("expected 120, got %v", out["duration"])
	}
}

func TestApply_ClampsAboveMax(t *testing.T) {
	c, _ := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "duration", Min: 0, Max: 100}},
	})
	out := c.Apply(base())
	if out["duration"] != float64(100) {
		t.Fatalf("expected 100, got %v", out["duration"])
	}
}

func TestApply_ClampsBelowMin(t *testing.T) {
	c, _ := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Min: 10, Max: 100}},
	})
	out := c.Apply(base())
	if out["score"] != float64(10) {
		t.Fatalf("expected 10, got %v", out["score"])
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	c, _ := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Min: 1, Max: 10}},
	})
	out := c.Apply(base())
	if out["score"] != float64(5) {
		t.Fatalf("expected 5, got %v", out["score"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	c, _ := clamp.New(clamp.Config{
		Rules:           []clamp.Rule{{Field: "DURATION", Min: 0, Max: 50}},
		CaseInsensitive: true,
	})
	out := c.Apply(base())
	if out["duration"] != float64(50) {
		t.Fatalf("expected 50, got %v", out["duration"])
	}
}

func TestApply_MissingField_NoChange(t *testing.T) {
	c, _ := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "latency", Min: 0, Max: 10}},
	})
	out := c.Apply(base())
	if _, ok := out["latency"]; ok {
		t.Fatal("expected latency to be absent")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	c, _ := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "duration", Min: 0, Max: 10}},
	})
	in := base()
	c.Apply(in)
	if in["duration"] != float64(120) {
		t.Fatal("original entry was mutated")
	}
}

func TestNew_InvalidRule_EmptyField(t *testing.T) {
	_, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "", Min: 0, Max: 10}},
	})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_InvalidRule_MinExceedsMax(t *testing.T) {
	_, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Min: 100, Max: 1}},
	})
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}
