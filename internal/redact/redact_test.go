package redact_test

import (
	"testing"

	"github.com/yourusername/logpipe/internal/reader"
	"github.com/yourusername/logpipe/internal/redact"
)

func baseEntry() reader.Entry {
	return reader.Entry{
		Timestamp: "2024-01-01T00:00:00Z",
		Level:     "info",
		Message:   "user login",
		Source:    "auth",
		Fields: map[string]any{
			"password": "s3cr3t",
			"token":    "abc123",
			"user":     "alice",
			"email":    "alice@example.com",
		},
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := redact.New(redact.Config{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestApply_ExactField_IsRedacted(t *testing.T) {
	r, _ := redact.New(redact.Config{Fields: []string{"password"}})
	out := r.Apply(baseEntry())
	if out.Fields["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out.Fields["password"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("non-redacted field altered: %v", out.Fields["user"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	r, _ := redact.New(redact.Config{Fields: []string{"PASSWORD"}})
	out := r.Apply(baseEntry())
	if out.Fields["password"] != "[REDACTED]" {
		t.Errorf("case-insensitive match failed: %v", out.Fields["password"])
	}
}

func TestApply_PatternMatch_IsRedacted(t *testing.T) {
	r, _ := redact.New(redact.Config{Patterns: []string{"^token$"}})
	out := r.Apply(baseEntry())
	if out.Fields["token"] != "[REDACTED]" {
		t.Errorf("pattern match failed: %v", out.Fields["token"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("non-matching field altered: %v", out.Fields["user"])
	}
}

func TestApply_CustomMask(t *testing.T) {
	r, _ := redact.New(redact.Config{Fields: []string{"email"}, Mask: "***"})
	out := r.Apply(baseEntry())
	if out.Fields["email"] != "***" {
		t.Errorf("custom mask not applied: %v", out.Fields["email"])
	}
}

func TestApply_NoConfig_PassesThrough(t *testing.T) {
	r, _ := redact.New(redact.Config{})
	e := baseEntry()
	out := r.Apply(e)
	if out.Fields["password"] != e.Fields["password"] {
		t.Errorf("field modified unexpectedly: %v", out.Fields["password"])
	}
}

func TestApply_OriginalEntryUnmodified(t *testing.T) {
	r, _ := redact.New(redact.Config{Fields: []string{"password"}})
	e := baseEntry()
	r.Apply(e)
	if e.Fields["password"] != "s3cr3t" {
		t.Error("original entry was mutated")
	}
}
