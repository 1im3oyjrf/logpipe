package remap_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/remap"
)

func base() map[string]any {
	return map[string]any{
		"level":   "WARN",
		"message": "disk pressure",
		"status":  "ERR",
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	r := remap.New(remap.Config{})
	out := r.Apply(base())
	if out["level"] != "WARN" {
		t.Fatalf("expected WARN, got %v", out["level"])
	}
}

func TestApply_MatchingValue_IsReplaced(t *testing.T) {
	r := remap.New(remap.Config{
		Rules: []remap.Rule{
			{Field: "status", From: "ERR", To: "error", CaseSensitive: true},
		},
	})
	out := r.Apply(base())
	if out["status"] != "error" {
		t.Fatalf("expected error, got %v", out["status"])
	}
}

func TestApply_NonMatchingValue_Unchanged(t *testing.T) {
	r := remap.New(remap.Config{
		Rules: []remap.Rule{
			{Field: "status", From: "OK", To: "ok", CaseSensitive: true},
		},
	})
	out := r.Apply(base())
	if out["status"] != "ERR" {
		t.Fatalf("expected ERR unchanged, got %v", out["status"])
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	r := remap.New(remap.Config{
		Rules: []remap.Rule{
			{Field: "level", From: "warn", To: "warning", CaseSensitive: false},
		},
	})
	out := r.Apply(base())
	if out["level"] != "warning" {
		t.Fatalf("expected warning, got %v", out["level"])
	}
}

func TestApply_MissingField_NoOp(t *testing.T) {
	r := remap.New(remap.Config{
		Rules: []remap.Rule{
			{Field: "nonexistent", From: "x", To: "y"},
		},
	})
	out := r.Apply(base())
	if _, ok := out["nonexistent"]; ok {
		t.Fatal("expected no nonexistent field")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	r := remap.New(remap.Config{
		Rules: []remap.Rule{
			{Field: "status", From: "ERR", To: "error", CaseSensitive: true},
		},
	})
	orig := base()
	r.Apply(orig)
	if orig["status"] != "ERR" {
		t.Fatal("original entry was mutated")
	}
}
