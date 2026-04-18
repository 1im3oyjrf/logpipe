package prefix_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/prefix"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello world",
		"service": "api",
	}
}

func TestApply_NoConfig_PassesThrough(t *testing.T) {
	p := prefix.New(prefix.Config{})
	entry := base()
	out := p.Apply(entry)
	if out["message"] != "hello world" {
		t.Fatalf("expected original message, got %v", out["message"])
	}
}

func TestApply_PrependsValue(t *testing.T) {
	p := prefix.New(prefix.Config{Field: "message", Value: "[APP]", Sep: " "})
	out := p.Apply(base())
	if out["message"] != "[APP] hello world" {
		t.Fatalf("unexpected value: %v", out["message"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	p := prefix.New(prefix.Config{Field: "MESSAGE", Value: ">>>", Sep: ""})
	out := p.Apply(base())
	if out["message"] != ">>>hello world" {
		t.Fatalf("unexpected value: %v", out["message"])
	}
}

func TestApply_MissingField_NoChange(t *testing.T) {
	p := prefix.New(prefix.Config{Field: "trace", Value: "[X]"})
	entry := base()
	out := p.Apply(entry)
	if _, ok := out["trace"]; ok {
		t.Fatal("trace field should not be injected")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	p := prefix.New(prefix.Config{Field: "service", Value: "env:", Sep: "-"})
	entry := base()
	p.Apply(entry)
	if entry["service"] != "api" {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_NonStringField_NoChange(t *testing.T) {
	p := prefix.New(prefix.Config{Field: "count", Value: "n="})
	entry := map[string]any{"count": 42}
	out := p.Apply(entry)
	if out["count"] != 42 {
		t.Fatalf("expected 42, got %v", out["count"])
	}
}
