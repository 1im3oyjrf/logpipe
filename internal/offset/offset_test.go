package offset_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/offset"
)

func base() map[string]any {
	return map[string]any{
		"level":    "info",
		"msg":      "hello",
		"duration": float64(100),
	}
}

func TestApply_NoField_PassesThrough(t *testing.T) {
	p := offset.New(offset.Config{By: 10})
	in := base()
	out := p.Apply(in)
	if out["duration"] != float64(100) {
		t.Fatalf("expected 100, got %v", out["duration"])
	}
}

func TestApply_ShiftsPositive(t *testing.T) {
	p := offset.New(offset.Config{Field: "duration", By: 50})
	out := p.Apply(base())
	if out["duration"] != float64(150) {
		t.Fatalf("expected 150, got %v", out["duration"])
	}
}

func TestApply_ShiftsNegative(t *testing.T) {
	p := offset.New(offset.Config{Field: "duration", By: -30})
	out := p.Apply(base())
	if out["duration"] != float64(70) {
		t.Fatalf("expected 70, got %v", out["duration"])
	}
}

func TestApply_MissingField_NoChange(t *testing.T) {
	p := offset.New(offset.Config{Field: "latency", By: 5})
	out := p.Apply(base())
	if _, ok := out["latency"]; ok {
		t.Fatal("expected no latency field")
	}
}

func TestApply_NonNumericField_NoChange(t *testing.T) {
	p := offset.New(offset.Config{Field: "msg", By: 1})
	out := p.Apply(base())
	if out["msg"] != "hello" {
		t.Fatalf("expected msg unchanged, got %v", out["msg"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	p := offset.New(offset.Config{Field: "DURATION", By: 10, CaseInsensitive: true})
	out := p.Apply(base())
	if out["duration"] != float64(110) {
		t.Fatalf("expected 110, got %v", out["duration"])
	}
}

func TestApply_CustomTarget_MovesField(t *testing.T) {
	p := offset.New(offset.Config{Field: "duration", By: 0, Target: "duration_ms"})
	out := p.Apply(base())
	if _, ok := out["duration"]; ok {
		t.Fatal("expected original field removed")
	}
	if out["duration_ms"] != float64(100) {
		t.Fatalf("expected 100, got %v", out["duration_ms"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	p := offset.New(offset.Config{Field: "duration", By: 1})
	in := base()
	p.Apply(in)
	if in["duration"] != float64(100) {
		t.Fatal("original entry was mutated")
	}
}
