package stencil_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/stencil"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
	}
}

func TestNew_EmptyFields_ReturnsError(t *testing.T) {
	_, err := stencil.New(stencil.Config{})
	if err == nil {
		t.Fatal("expected error for empty Fields, got nil")
	}
}

func TestApply_InjectsNewFields(t *testing.T) {
	s, err := stencil.New(stencil.Config{
		Fields: map[string]any{"env": "prod", "region": "us-east-1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Apply(base())
	if out["env"] != "prod" {
		t.Errorf("env: got %v, want prod", out["env"])
	}
	if out["region"] != "us-east-1" {
		t.Errorf("region: got %v, want us-east-1", out["region"])
	}
}

func TestApply_DoesNotOverwriteByDefault(t *testing.T) {
	s, _ := stencil.New(stencil.Config{
		Fields: map[string]any{"level": "debug"},
	})
	out := s.Apply(base())
	if out["level"] != "info" {
		t.Errorf("level should not be overwritten: got %v", out["level"])
	}
}

func TestApply_OverwriteReplacesExisting(t *testing.T) {
	s, _ := stencil.New(stencil.Config{
		Fields:    map[string]any{"level": "warn"},
		Overwrite: true,
	})
	out := s.Apply(base())
	if out["level"] != "warn" {
		t.Errorf("level: got %v, want warn", out["level"])
	}
}

func TestApply_CaseInsensitive_MatchesExistingKey(t *testing.T) {
	s, _ := stencil.New(stencil.Config{
		Fields:          map[string]any{"LEVEL": "error"},
		Overwrite:       true,
		CaseInsensitive: true,
	})
	out := s.Apply(base())
	if out["level"] != "error" {
		t.Errorf("level: got %v, want error", out["level"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	s, _ := stencil.New(stencil.Config{
		Fields: map[string]any{"env": "staging"},
	})
	orig := base()
	s.Apply(orig)
	if _, ok := orig["env"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_NoExistingFields_AllInjected(t *testing.T) {
	s, _ := stencil.New(stencil.Config{
		Fields: map[string]any{"a": 1, "b": 2},
	})
	out := s.Apply(map[string]any{})
	if out["a"] != 1 || out["b"] != 2 {
		t.Errorf("expected a=1 b=2, got %v", out)
	}
}
