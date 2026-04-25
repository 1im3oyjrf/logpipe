package fold_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/fold"
)

var base = map[string]any{
	"first": "hello",
	"last":  "world",
	"level": "info",
}

func TestNew_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := fold.New(fold.Config{Fields: []string{"a", "b"}, Target: ""})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestNew_TooFewFields_ReturnsError(t *testing.T) {
	_, err := fold.New(fold.Config{Fields: []string{"a"}, Target: "out"})
	if err == nil {
		t.Fatal("expected error for fewer than two source fields")
	}
}

func TestApply_FoldsFieldsIntoTarget(t *testing.T) {
	f, err := fold.New(fold.Config{Fields: []string{"first", "last"}, Target: "full"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := f.Apply(base)
	if got := out["full"]; got != "hello world" {
		t.Errorf("want %q, got %q", "hello world", got)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f, _ := fold.New(fold.Config{Fields: []string{"first", "last"}, Target: "full", Separator: "-"})
	out := f.Apply(base)
	if got := out["full"]; got != "hello-world" {
		t.Errorf("want %q, got %q", "hello-world", got)
	}
}

func TestApply_DropSources_RemovesFields(t *testing.T) {
	f, _ := fold.New(fold.Config{Fields: []string{"first", "last"}, Target: "full", DropSources: true})
	out := f.Apply(base)
	if _, ok := out["first"]; ok {
		t.Error("expected 'first' to be removed")
	}
	if _, ok := out["last"]; ok {
		t.Error("expected 'last' to be removed")
	}
	if out["full"] != "hello world" {
		t.Errorf("unexpected full value: %v", out["full"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	entry := map[string]any{"First": "hello", "Last": "world"}
	f, _ := fold.New(fold.Config{Fields: []string{"first", "last"}, Target: "full", CaseInsensitive: true})
	out := f.Apply(entry)
	if got := out["full"]; got != "hello world" {
		t.Errorf("want %q, got %q", "hello world", got)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	f, _ := fold.New(fold.Config{Fields: []string{"first", "last"}, Target: "full", DropSources: true})
	_ = f.Apply(base)
	if _, ok := base["first"]; !ok {
		t.Error("original entry was mutated")
	}
}

func TestApply_MissingFields_NoOutput(t *testing.T) {
	f, _ := fold.New(fold.Config{Fields: []string{"x", "y"}, Target: "full"})
	out := f.Apply(base)
	if _, ok := out["full"]; ok {
		t.Error("expected no target field when sources are missing")
	}
}
