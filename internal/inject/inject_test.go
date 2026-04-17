package inject_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/inject"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello world",
	}
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	inj := inject.New(inject.Config{})
	entry := base()
	out := inj.Apply(entry)
	if len(out) != len(entry) {
		t.Fatalf("expected %d fields, got %d", len(entry), len(out))
	}
}

func TestApply_InjectsStaticField(t *testing.T) {
	inj := inject.New(inject.Config{
		Fields: map[string]string{"env": "production"},
	})
	out := inj.Apply(base())
	if out["env"] != "production" {
		t.Fatalf("expected env=production, got %v", out["env"])
	}
}

func TestApply_CaseInsensitiveKey(t *testing.T) {
	inj := inject.New(inject.Config{
		Fields: map[string]string{"ENV": "staging"},
	})
	out := inj.Apply(base())
	if out["env"] != "staging" {
		t.Fatalf("expected env=staging, got %v", out["env"])
	}
}

func TestApply_DoesNotOverwriteByDefault(t *testing.T) {
	inj := inject.New(inject.Config{
		Fields: map[string]string{"level": "debug"},
	})
	out := inj.Apply(base())
	if out["level"] != "info" {
		t.Fatalf("expected level=info (not overwritten), got %v", out["level"])
	}
}

func TestApply_OverwriteExisting(t *testing.T) {
	inj := inject.New(inject.Config{
		Fields:            map[string]string{"level": "debug"},
		OverwriteExisting: true,
	})
	out := inj.Apply(base())
	if out["level"] != "debug" {
		t.Fatalf("expected level=debug, got %v", out["level"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	inj := inject.New(inject.Config{
		Fields: map[string]string{"env": "test"},
	})
	entry := base()
	inj.Apply(entry)
	if _, ok := entry["env"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_MultipleFields(t *testing.T) {
	inj := inject.New(inject.Config{
		Fields: map[string]string{"env": "prod", "region": "us-east-1"},
	})
	out := inj.Apply(base())
	if out["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", out["env"])
	}
	if out["region"] != "us-east-1" {
		t.Fatalf("expected region=us-east-1, got %v", out["region"])
	}
}
