package coalesce_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/coalesce"
)

func base() map[string]any {
	return map[string]any{
		"message": "hello",
		"level":   "info",
	}
}

func TestApply_NoRules_PassesThrough(t *testing.T) {
	tr := coalesce.New(coalesce.Config{})
	in := base()
	out := tr.Apply(in)
	if out["message"] != "hello" {
		t.Fatalf("expected message=hello, got %v", out["message"])
	}
}

func TestApply_FirstSourceUsed(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"msg", "message"}, Target: "canonical"},
		},
	})
	in := map[string]any{"msg": "first", "message": "second"}
	out := tr.Apply(in)
	if out["canonical"] != "first" {
		t.Fatalf("expected canonical=first, got %v", out["canonical"])
	}
}

func TestApply_FallsBackToSecondSource(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"msg", "message"}, Target: "canonical"},
		},
	})
	in := map[string]any{"message": "fallback"}
	out := tr.Apply(in)
	if out["canonical"] != "fallback" {
		t.Fatalf("expected canonical=fallback, got %v", out["canonical"])
	}
}

func TestApply_SourcesRemovedByDefault(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"msg", "message"}, Target: "canonical", KeepSources: false},
		},
	})
	in := map[string]any{"msg": "hi", "message": "there"}
	out := tr.Apply(in)
	if _, ok := out["msg"]; ok {
		t.Fatal("expected msg to be removed")
	}
}

func TestApply_KeepSources_RetainsThem(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"msg", "message"}, Target: "canonical", KeepSources: true},
		},
	})
	in := map[string]any{"msg": "hi", "message": "there"}
	out := tr.Apply(in)
	if out["msg"] != "hi" {
		t.Fatal("expected msg to be kept")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"msg"}, Target: "canonical"},
		},
	})
	in := map[string]any{"msg": "hi"}
	_ = tr.Apply(in)
	if _, ok := in["canonical"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_AllSourcesEmpty_TargetNotSet(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"a", "b"}, Target: "canonical"},
		},
	})
	in := map[string]any{"level": "info"}
	out := tr.Apply(in)
	if _, ok := out["canonical"]; ok {
		t.Fatal("expected canonical not to be set")
	}
}

// TestApply_TargetSameAsSource verifies that when the target field name matches
// one of the sources, the value is correctly promoted and the source is removed.
func TestApply_TargetSameAsSource(t *testing.T) {
	tr := coalesce.New(coalesce.Config{
		Rules: []coalesce.Rule{
			{Sources: []string{"msg", "message"}, Target: "message", KeepSources: false},
		},
	})
	in := map[string]any{"msg": "preferred"}
	out := tr.Apply(in)
	if out["message"] != "preferred" {
		t.Fatalf("expected message=preferred, got %v", out["message"])
	}
	if _, ok := out["msg"]; ok {
		t.Fatal("expected msg to be removed")
	}
}
