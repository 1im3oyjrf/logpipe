package mask_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/mask"
)

func base() map[string]any {
	return map[string]any{
		"level":    "info",
		"message":  "user login",
		"password": "s3cr3t",
		"token":    "abc123",
	}
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	m := mask.New(mask.Config{})
	in := base()
	out := m.Apply(in)
	if out["password"] != "s3cr3t" {
		t.Errorf("expected original value, got %v", out["password"])
	}
}

func TestApply_MasksConfiguredField(t *testing.T) {
	m := mask.New(mask.Config{Fields: []string{"password"}})
	out := m.Apply(base())
	if out["password"] != "***" {
		t.Errorf("expected ***, got %v", out["password"])
	}
	if out["token"] != "abc123" {
		t.Errorf("token should be unmasked")
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	m := mask.New(mask.Config{Fields: []string{"PASSWORD"}})
	out := m.Apply(base())
	if out["password"] != "***" {
		t.Errorf("expected case-insensitive mask, got %v", out["password"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	m := mask.New(mask.Config{Fields: []string{"token"}, Placeholder: "[REDACTED]"})
	out := m.Apply(base())
	if out["token"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out["token"])
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	m := mask.New(mask.Config{Fields: []string{"password", "token"}})
	in := base()
	_ = m.Apply(in)
	if in["password"] != "s3cr3t" {
		t.Error("original entry was mutated")
	}
}

func TestApply_MissingField_Ignored(t *testing.T) {
	m := mask.New(mask.Config{Fields: []string{"nonexistent"}})
	out := m.Apply(base())
	if _, ok := out["nonexistent"]; ok {
		t.Error("nonexistent field should not appear in output")
	}
	if out["level"] != "info" {
		t.Error("existing fields should be preserved")
	}
}
