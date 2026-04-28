package quantize_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/quantize"
)

func base() map[string]any {
	return map[string]any{
		"level":    "info",
		"message":  "request handled",
		"duration": 137.8,
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := quantize.New(quantize.Config{Step: 10})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NegativeStep_ReturnsError(t *testing.T) {
	_, err := quantize.New(quantize.Config{Field: "duration", Step: -5})
	if err == nil {
		t.Fatal("expected error for negative step")
	}
}

func TestNew_ZeroStep_DefaultsToOne(t *testing.T) {
	q, err := quantize.New(quantize.Config{Field: "duration"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := q.Apply(base())
	if got, ok := out["duration"].(float64); !ok || got != 137.0 {
		t.Fatalf("expected 137.0, got %v", out["duration"])
	}
}

func TestApply_QuantizesToStep(t *testing.T) {
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 50})
	out := q.Apply(base())
	if got := out["duration"].(float64); got != 100.0 {
		t.Fatalf("expected 100.0, got %v", got)
	}
}

func TestApply_ExactMultiple_Unchanged(t *testing.T) {
	e := map[string]any{"duration": 200.0}
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 100})
	out := q.Apply(e)
	if got := out["duration"].(float64); got != 200.0 {
		t.Fatalf("expected 200.0, got %v", got)
	}
}

func TestApply_CustomTarget_WritesToTarget(t *testing.T) {
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 100, Target: "duration_bucket"})
	out := q.Apply(base())
	if _, ok := out["duration_bucket"]; !ok {
		t.Fatal("expected duration_bucket field")
	}
	if _, ok := out["duration"]; !ok {
		t.Fatal("original duration field should be preserved")
	}
}

func TestApply_DoesNotOverwriteByDefault(t *testing.T) {
	e := map[string]any{"duration": 137.8, "bucket": 999.0}
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 50, Target: "bucket"})
	out := q.Apply(e)
	if got := out["bucket"].(float64); got != 999.0 {
		t.Fatalf("expected original bucket 999.0, got %v", got)
	}
}

func TestApply_Overwrite_ReplacesExisting(t *testing.T) {
	e := map[string]any{"duration": 137.8, "bucket": 999.0}
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 50, Target: "bucket", Overwrite: true})
	out := q.Apply(e)
	if got := out["bucket"].(float64); got != 100.0 {
		t.Fatalf("expected 100.0, got %v", got)
	}
}

func TestApply_MissingField_PassesThrough(t *testing.T) {
	e := map[string]any{"level": "info"}
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 50})
	out := q.Apply(e)
	if _, ok := out["duration"]; ok {
		t.Fatal("duration should not be injected when missing")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	orig := base()
	q, _ := quantize.New(quantize.Config{Field: "duration", Step: 50})
	_ = q.Apply(orig)
	if orig["duration"].(float64) != 137.8 {
		t.Fatal("original entry was mutated")
	}
}
