package prune_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/prune"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
		"secret":  "topsecret",
		"token":   "abc123",
	}
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	p := prune.New(prune.Config{})
	out := p.Apply(base())
	if len(out) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(out))
	}
}

func TestApply_RemovesConfiguredField(t *testing.T) {
	p := prune.New(prune.Config{Fields: []string{"secret"}})
	out := p.Apply(base())
	if _, ok := out["secret"]; ok {
		t.Fatal("expected secret to be removed")
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(out))
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	p := prune.New(prune.Config{Fields: []string{"TOKEN"}})
	out := p.Apply(base())
	if _, ok := out["token"]; ok {
		t.Fatal("expected token to be removed")
	}
}

func TestApply_MultipleFields(t *testing.T) {
	p := prune.New(prune.Config{Fields: []string{"secret", "token"}})
	out := p.Apply(base())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	orig := base()
	p := prune.New(prune.Config{Fields: []string{"secret"}})
	_ = p.Apply(orig)
	if _, ok := orig["secret"]; !ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestApply_UnknownField_Ignored(t *testing.T) {
	p := prune.New(prune.Config{Fields: []string{"nonexistent"}})
	out := p.Apply(base())
	if len(out) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(out))
	}
}
