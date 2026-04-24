package clamp_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/clamp"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.Entry {
	return reader.Entry{
		"level":   "info",
		"message": "test",
		"score":   42.0,
		"count":   100.0,
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	c, err := clamp.New(clamp.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	out := c.Apply(e)
	if out["score"] != e["score"] {
		t.Errorf("expected score unchanged, got %v", out["score"])
	}
}

func TestApply_ClampsAboveMax(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Max: ptr(50.0)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	e["score"] = 99.0
	out := c.Apply(e)
	if out["score"] != 50.0 {
		t.Errorf("expected 50.0, got %v", out["score"])
	}
}

func TestApply_ClampsBelowMin(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Min: ptr(10.0)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	e["score"] = 1.0
	out := c.Apply(e)
	if out["score"] != 10.0 {
		t.Errorf("expected 10.0, got %v", out["score"])
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Min: ptr(0.0), Max: ptr(100.0)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	e["score"] = 42.0
	out := c.Apply(e)
	if out["score"] != 42.0 {
		t.Errorf("expected 42.0, got %v", out["score"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "SCORE", Max: ptr(50.0)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	e["score"] = 99.0
	out := c.Apply(e)
	if out["score"] != 50.0 {
		t.Errorf("expected clamped value 50.0, got %v", out["score"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "score", Max: ptr(50.0)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	e["score"] = 99.0
	_ = c.Apply(e)
	if e["score"] != 99.0 {
		t.Errorf("original entry was mutated")
	}
}

func TestApply_MultipleRules_AllApplied(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{
			{Field: "score", Max: ptr(50.0)},
			{Field: "count", Min: ptr(200.0)},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	e["score"] = 99.0
	e["count"] = 100.0
	out := c.Apply(e)
	if out["score"] != 50.0 {
		t.Errorf("expected score 50.0, got %v", out["score"])
	}
	if out["count"] != 200.0 {
		t.Errorf("expected count 200.0, got %v", out["count"])
	}
}

func TestApply_NonNumericField_PassesThrough(t *testing.T) {
	c, err := clamp.New(clamp.Config{
		Rules: []clamp.Rule{{Field: "message", Max: ptr(10.0)}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := base()
	out := c.Apply(e)
	if out["message"] != "test" {
		t.Errorf("expected message unchanged, got %v", out["message"])
	}
}

func ptr(f float64) *float64 { return &f }
