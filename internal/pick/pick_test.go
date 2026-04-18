package pick_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/pick"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
		"service": "api",
		"host":    "box1",
	}
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	p := pick.New(pick.Config{})
	out := p.Apply(base())
	if len(out) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(out))
	}
}

func TestApply_KeepsOnlyConfiguredFields(t *testing.T) {
	p := pick.New(pick.Config{Fields: []string{"level", "message"}})
	out := p.Apply(base())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if _, ok := out["level"]; !ok {
		t.Error("expected 'level' to be present")
	}
	if _, ok := out["message"]; !ok {
		t.Error("expected 'message' to be present")
	}
	if _, ok := out["service"]; ok {
		t.Error("expected 'service' to be absent")
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	p := pick.New(pick.Config{Fields: []string{"LEVEL", "Service"}})
	out := p.Apply(base())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	orig := base()
	p := pick.New(pick.Config{Fields: []string{"level"}})
	_ = p.Apply(orig)
	if len(orig) != 4 {
		t.Error("original entry was mutated")
	}
}

func TestApply_MissingField_NotInOutput(t *testing.T) {
	p := pick.New(pick.Config{Fields: []string{"level", "nonexistent"}})
	out := p.Apply(base())
	if _, ok := out["nonexistent"]; ok {
		t.Error("did not expect 'nonexistent' in output")
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 field, got %d", len(out))
	}
}
