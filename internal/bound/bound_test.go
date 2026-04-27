package bound_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/bound"
)

func base() map[string]any {
	return map[string]any{"level": "info", "latency": 42.0}
}

func ptr(f float64) *float64 { return &f }

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := bound.New(bound.Config{Min: ptr(0)})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NoBounds_ReturnsError(t *testing.T) {
	_, err := bound.New(bound.Config{Field: "latency"})
	if err == nil {
		t.Fatal("expected error when neither Min nor Max is set")
	}
}

func TestNew_MinExceedsMax_ReturnsError(t *testing.T) {
	_, err := bound.New(bound.Config{Field: "latency", Min: ptr(10), Max: ptr(5)})
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "latency", Min: ptr(0), Max: ptr(100)})
	out := b.Apply(base())
	if out["latency"] != 42.0 {
		t.Fatalf("expected 42.0, got %v", out["latency"])
	}
}

func TestApply_AboveMax_IsClamped(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "latency", Max: ptr(30)})
	out := b.Apply(base())
	if out["latency"] != 30.0 {
		t.Fatalf("expected 30.0, got %v", out["latency"])
	}
}

func TestApply_BelowMin_IsClamped(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "latency", Min: ptr(50)})
	out := b.Apply(base())
	if out["latency"] != 50.0 {
		t.Fatalf("expected 50.0, got %v", out["latency"])
	}
}

func TestApply_FlagField_SetOnClamp(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "latency", Max: ptr(10), FlagField: "clamped"})
	out := b.Apply(base())
	if out["clamped"] != "true" {
		t.Fatalf("expected clamped=true, got %v", out["clamped"])
	}
}

func TestApply_FlagField_FalseWhenNotClamped(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "latency", Min: ptr(0), Max: ptr(100), FlagField: "clamped"})
	out := b.Apply(base())
	if out["clamped"] != "false" {
		t.Fatalf("expected clamped=false, got %v", out["clamped"])
	}
}

func TestApply_MissingField_PassesThrough(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "missing", Min: ptr(0)})
	in := base()
	out := b.Apply(in)
	if _, ok := out["missing"]; ok {
		t.Fatal("unexpected field injection for missing key")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	b, _ := bound.New(bound.Config{Field: "latency", Max: ptr(10)})
	in := base()
	b.Apply(in)
	if in["latency"] != 42.0 {
		t.Fatal("original entry was mutated")
	}
}
