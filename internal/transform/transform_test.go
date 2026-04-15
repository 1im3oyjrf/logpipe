package transform_test

import (
	"testing"
	"time"

	"github.com/yourusername/logpipe/internal/reader"
	"github.com/yourusername/logpipe/internal/transform"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "hello world",
		Source:    "app",
		Fields: map[string]any{
			"user":    "alice",
			"token":   "secret-abc",
			"request": "/api/v1/data",
		},
	}
}

func TestApply_NoConfig_PassesThrough(t *testing.T) {
	tr := transform.New(transform.Config{})
	e := baseEntry()
	out := tr.Apply(e)

	if out.Message != e.Message {
		t.Errorf("expected message %q, got %q", e.Message, out.Message)
	}
	if len(out.Fields) != len(e.Fields) {
		t.Errorf("expected %d fields, got %d", len(e.Fields), len(out.Fields))
	}
}

func TestApply_RedactField(t *testing.T) {
	tr := transform.New(transform.Config{
		RedactFields: []string{"token"},
	})
	out := tr.Apply(baseEntry())

	if out.Fields["token"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out.Fields["token"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("non-redacted field should be unchanged")
	}
}

func TestApply_RedactField_CaseInsensitive(t *testing.T) {
	tr := transform.New(transform.Config{
		RedactFields: []string{"TOKEN"},
	})
	out := tr.Apply(baseEntry())

	if out.Fields["token"] != "[REDACTED]" {
		t.Errorf("redaction should be case-insensitive")
	}
}

func TestApply_RenameField(t *testing.T) {
	tr := transform.New(transform.Config{
		RenameFields: map[string]string{"user": "username"},
	})
	out := tr.Apply(baseEntry())

	if _, ok := out.Fields["user"]; ok {
		t.Error("old field name should not exist after rename")
	}
	if out.Fields["username"] != "alice" {
		t.Errorf("expected renamed field value %q, got %v", "alice", out.Fields["username"])
	}
}

func TestApply_AddFields(t *testing.T) {
	tr := transform.New(transform.Config{
		AddFields: map[string]string{"env": "production"},
	})
	out := tr.Apply(baseEntry())

	if out.Fields["env"] != "production" {
		t.Errorf("expected injected field 'env'='production', got %v", out.Fields["env"])
	}
}

func TestApply_AddFields_DoesNotOverwrite(t *testing.T) {
	tr := transform.New(transform.Config{
		AddFields: map[string]string{"user": "injected"},
	})
	out := tr.Apply(baseEntry())

	if out.Fields["user"] != "alice" {
		t.Errorf("AddFields should not overwrite existing field; got %v", out.Fields["user"])
	}
}

func TestApply_OriginalEntryUnmodified(t *testing.T) {
	tr := transform.New(transform.Config{
		RedactFields: []string{"token"},
		RenameFields: map[string]string{"user": "username"},
	})
	e := baseEntry()
	_ = tr.Apply(e)

	if e.Fields["token"] == "[REDACTED]" {
		t.Error("original entry should not be mutated by Apply")
	}
	if _, ok := e.Fields["user"]; !ok {
		t.Error("original entry field 'user' should still exist")
	}
}
