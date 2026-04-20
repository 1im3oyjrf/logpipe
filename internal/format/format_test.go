package format_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/format"
)

var base = map[string]any{
	"level":   "error",
	"message": "disk full",
	"host":    "web-01",
}

func TestNew_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := format.New(format.Config{Target: "", Template: "{level}: {message}"})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestNew_EmptyTemplate_ReturnsError(t *testing.T) {
	_, err := format.New(format.Config{Target: "summary", Template: ""})
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestApply_RendersTemplate(t *testing.T) {
	f, err := format.New(format.Config{Target: "summary", Template: "[{level}] {message}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := f.Apply(base)
	want := "[error] disk full"
	if got := out["summary"]; got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	f, _ := format.New(format.Config{Target: "summary", Template: "{host}"})
	origLen := len(base)
	f.Apply(base)
	if len(base) != origLen {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_DoesNotOverwriteByDefault(t *testing.T) {
	entry := map[string]any{"level": "info", "summary": "existing"}
	f, _ := format.New(format.Config{Target: "summary", Template: "[{level}]"})
	out := f.Apply(entry)
	if got := out["summary"]; got != "existing" {
		t.Fatalf("expected existing value preserved, got %q", got)
	}
}

func TestApply_OverwriteReplaces(t *testing.T) {
	entry := map[string]any{"level": "warn", "summary": "old"}
	f, _ := format.New(format.Config{Target: "summary", Template: "[{level}]", Overwrite: true})
	out := f.Apply(entry)
	if got := out["summary"]; got != "[warn]" {
		t.Fatalf("want \"[warn]\", got %q", got)
	}
}

func TestApply_MissingPlaceholder_LeftBlank(t *testing.T) {
	f, _ := format.New(format.Config{Target: "out", Template: "host={host} svc={service}"})
	out := f.Apply(map[string]any{"host": "srv1"})
	want := "host=srv1 svc={service}"
	if got := out["out"]; got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
