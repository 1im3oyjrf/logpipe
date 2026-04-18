package merge_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/merge"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
	}
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	m := merge.New(merge.Config{})
	out := m.Apply(base())
	if out["level"] != "info" {
		t.Fatalf("expected info, got %v", out["level"])
	}
}

func TestApply_InjectsNewField(t *testing.T) {
	m := merge.New(merge.Config{Fields: map[string]string{"env": "prod"}})
	out := m.Apply(base())
	if out["env"] != "prod" {
		t.Fatalf("expected prod, got %v", out["env"])
	}
}

func TestApply_DoesNotOverwriteByDefault(t *testing.T) {
	m := merge.New(merge.Config{Fields: map[string]string{"level": "debug"}})
	out := m.Apply(base())
	if out["level"] != "info" {
		t.Fatalf("expected original info, got %v", out["level"])
	}
}

func TestApply_OverwriteReplaces(t *testing.T) {
	m := merge.New(merge.Config{Fields: map[string]string{"level": "debug"}, Overwrite: true})
	out := m.Apply(base())
	if out["level"] != "debug" {
		t.Fatalf("expected debug, got %v", out["level"])
	}
}

func TestApply_CaseInsensitiveKey(t *testing.T) {
	m := merge.New(merge.Config{Fields: map[string]string{"ENV": "staging"}})
	out := m.Apply(base())
	if out["env"] != "staging" {
		t.Fatalf("expected staging under normalised key, got %v", out["env"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	entry := base()
	m := merge.New(merge.Config{Fields: map[string]string{"env": "prod"}})
	_ = m.Apply(entry)
	if _, ok := entry["env"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}
