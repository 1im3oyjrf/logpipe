package jitter_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/jitter"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
	}
}

func TestApply_DefaultField_InjectsJitterMS(t *testing.T) {
	a := jitter.New(jitter.Config{})
	out := a.Apply(base())
	v, ok := out["jitter_ms"]
	if !ok {
		t.Fatal("expected jitter_ms field to be present")
	}
	ms, ok := v.(int)
	if !ok {
		t.Fatalf("expected int, got %T", v)
	}
	if ms < 0 || ms >= 100 {
		t.Errorf("value %d out of expected range [0,100)", ms)
	}
}

func TestApply_CustomField_UsesField(t *testing.T) {
	a := jitter.New(jitter.Config{Field: "delay", MaxMS: 50})
	out := a.Apply(base())
	if _, ok := out["delay"]; !ok {
		t.Fatal("expected delay field")
	}
	if _, ok := out["jitter_ms"]; ok {
		t.Fatal("unexpected jitter_ms field")
	}
}

func TestApply_MaxMS_BoundsRespected(t *testing.T) {
	a := jitter.New(jitter.Config{MaxMS: 10})
	for i := 0; i < 200; i++ {
		out := a.Apply(base())
		ms := out["jitter_ms"].(int)
		if ms < 0 || ms >= 10 {
			t.Fatalf("value %d out of range [0,10)", ms)
		}
	}
}

func TestApply_ExistingField_NotOverwrittenByDefault(t *testing.T) {
	a := jitter.New(jitter.Config{})
	entry := base()
	entry["jitter_ms"] = 999
	out := a.Apply(entry)
	if out["jitter_ms"] != 999 {
		t.Errorf("expected existing value to be preserved, got %v", out["jitter_ms"])
	}
}

func TestApply_Overwrite_ReplacesExisting(t *testing.T) {
	a := jitter.New(jitter.Config{Overwrite: true, MaxMS: 5})
	entry := base()
	entry["jitter_ms"] = 999
	out := a.Apply(entry)
	if out["jitter_ms"] == 999 {
		t.Error("expected existing value to be replaced")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	a := jitter.New(jitter.Config{})
	original := base()
	a.Apply(original)
	if _, ok := original["jitter_ms"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestNew_NegativeMaxMS_UsesDefault(t *testing.T) {
	a := jitter.New(jitter.Config{MaxMS: -5})
	out := a.Apply(base())
	ms := out["jitter_ms"].(int)
	if ms < 0 || ms >= 100 {
		t.Errorf("expected default range [0,100), got %d", ms)
	}
}
