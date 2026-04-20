package format_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/format"
)

func TestNew_ExplicitTarget_Stored(t *testing.T) {
	f, err := format.New(format.Config{Target: "label", Template: "{level}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := format.TargetOf(f); got != "label" {
		t.Fatalf("want \"label\", got %q", got)
	}
}

func TestNew_ExplicitTemplate_Stored(t *testing.T) {
	tmpl := "{host}/{level}"
	f, err := format.New(format.Config{Target: "path", Template: tmpl})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := format.TemplateOf(f); got != tmpl {
		t.Fatalf("want %q, got %q", tmpl, got)
	}
}

func TestNew_WhitespaceTarget_ReturnsError(t *testing.T) {
	_, err := format.New(format.Config{Target: "   ", Template: "{level}"})
	if err == nil {
		t.Fatal("expected error for whitespace-only target")
	}
}
