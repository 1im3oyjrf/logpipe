package typecast_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/typecast"
)

func base() map[string]any {
	return map[string]any{
		"message": "hello",
		"count":   "42",
		"ratio":   "3.14",
		"enabled": "true",
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	c := typecast.New(typecast.Config{})
	in := base()
	out := c.Apply(in)
	if out["count"] != "42" {
		t.Fatalf("expected string \"42\", got %v", out["count"])
	}
}

func TestApply_CoercesToInt(t *testing.T) {
	c := typecast.New(typecast.Config{Rules: []typecast.Rule{{Field: "count", Target: "int"}}})
	out := c.Apply(base())
	v, ok := out["count"].(int64)
	if !ok || v != 42 {
		t.Fatalf("expected int64(42), got %v (%T)", out["count"], out["count"])
	}
}

func TestApply_CoercesToFloat(t *testing.T) {
	c := typecast.New(typecast.Config{Rules: []typecast.Rule{{Field: "ratio", Target: "float"}}})
	out := c.Apply(base())
	v, ok := out["ratio"].(float64)
	if !ok || v != 3.14 {
		t.Fatalf("expected float64(3.14), got %v", out["ratio"])
	}
}

func TestApply_CoercesToBool(t *testing.T) {
	c := typecast.New(typecast.Config{Rules: []typecast.Rule{{Field: "enabled", Target: "bool"}}})
	out := c.Apply(base())
	v, ok := out["enabled"].(bool)
	if !ok || !v {
		t.Fatalf("expected bool(true), got %v", out["enabled"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	c := typecast.New(typecast.Config{Rules: []typecast.Rule{{Field: "COUNT", Target: "int"}}})
	out := c.Apply(base())
	v, ok := out["count"].(int64)
	if !ok || v != 42 {
		t.Fatalf("expected int64(42), got %v", out["count"])
	}
}

func TestApply_InvalidValue_SkipsField(t *testing.T) {
	in := map[string]any{"count": "not-a-number"}
	c := typecast.New(typecast.Config{Rules: []typecast.Rule{{Field: "count", Target: "int"}}})
	out := c.Apply(in)
	if out["count"] != "not-a-number" {
		t.Fatalf("expected original value preserved, got %v", out["count"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	in := base()
	c := typecast.New(typecast.Config{Rules: []typecast.Rule{{Field: "count", Target: "int"}}})
	c.Apply(in)
	if in["count"] != "42" {
		t.Fatal("original entry was mutated")
	}
}
