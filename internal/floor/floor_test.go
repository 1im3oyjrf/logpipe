package floor_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/floor"
)

func base() map[string]any {
	return map[string]any{
		"level":    "info",
		"latency":  42.0,
		"retries":  0.0,
		"priority": 3.0,
	}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "", Min: 1}},
	})
	if err == nil {
		t.Fatal("expected error for empty field name, got nil")
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	f, _ := floor.New(floor.Config{})
	in := base()
	out := f.Apply(in)
	if out["latency"] != in["latency"] {
		t.Fatalf("expected latency unchanged, got %v", out["latency"])
	}
}

func TestApply_BelowMin_IsClamped(t *testing.T) {
	f, _ := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "retries", Min: 1}},
	})
	out := f.Apply(base())
	if out["retries"] != 1.0 {
		t.Fatalf("expected retries=1.0, got %v", out["retries"])
	}
}

func TestApply_AboveMin_Unchanged(t *testing.T) {
	f, _ := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "latency", Min: 10}},
	})
	out := f.Apply(base())
	if out["latency"] != 42.0 {
		t.Fatalf("expected latency=42.0, got %v", out["latency"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	f, _ := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "PRIORITY", Min: 5}},
	})
	out := f.Apply(base())
	if out["priority"] != 5.0 {
		t.Fatalf("expected priority=5.0, got %v", out["priority"])
	}
}

func TestApply_MissingField_NoChange(t *testing.T) {
	f, _ := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "missing", Min: 10}},
	})
	out := f.Apply(base())
	if _, ok := out["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	f, _ := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "retries", Min: 5}},
	})
	in := base()
	f.Apply(in)
	if in["retries"] != 0.0 {
		t.Fatalf("original entry was mutated: retries=%v", in["retries"])
	}
}

func TestApply_NonNumericField_Skipped(t *testing.T) {
	f, _ := floor.New(floor.Config{
		Rules: []floor.Rule{{Field: "level", Min: 1}},
	})
	out := f.Apply(base())
	if out["level"] != "info" {
		t.Fatalf("expected level=info unchanged, got %v", out["level"])
	}
}
