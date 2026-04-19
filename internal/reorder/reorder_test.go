package reorder_test

import (
	"testing"

	"logpipe/internal/reader"
	"logpipe/internal/reorder"
)

func base() reader.Entry {
	return reader.Entry{
		Level:   "info",
		Message: "hello",
		Fields: map[string]any{
			"alpha": "a",
			"beta":  "b",
			"gamma": "g",
		},
	}
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	r := reorder.New(reorder.Config{})
	e := base()
	out := r.Apply(e)
	if len(out.Fields) != len(e.Fields) {
		t.Fatalf("expected %d fields, got %d", len(e.Fields), len(out.Fields))
	}
}

func TestApply_PromotesConfiguredFields(t *testing.T) {
	r := reorder.New(reorder.Config{Fields: []string{"gamma", "alpha"}})
	out := r.Apply(base())
	if out.Fields["gamma"] != "g" {
		t.Errorf("expected gamma=g, got %v", out.Fields["gamma"])
	}
	if out.Fields["alpha"] != "a" {
		t.Errorf("expected alpha=a, got %v", out.Fields["alpha"])
	}
	if out.Fields["beta"] != "b" {
		t.Errorf("expected beta=b, got %v", out.Fields["beta"])
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	r := reorder.New(reorder.Config{
		Fields:          []string{"GAMMA"},
		CaseInsensitive: true,
	})
	out := r.Apply(base())
	if out.Fields["gamma"] != "g" {
		t.Errorf("expected gamma promoted, got %v", out.Fields["gamma"])
	}
}

func TestApply_CaseSensitive_NoMatch(t *testing.T) {
	r := reorder.New(reorder.Config{
		Fields:          []string{"GAMMA"},
		CaseInsensitive: false,
	})
	out := r.Apply(base())
	// All fields still present, none promoted under wrong case.
	if len(out.Fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(out.Fields))
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	r := reorder.New(reorder.Config{Fields: []string{"beta"}})
	e := base()
	r.Apply(e)
	if e.Fields["alpha"] != "a" {
		t.Error("original entry was mutated")
	}
}

func TestApply_MissingConfiguredField_Ignored(t *testing.T) {
	r := reorder.New(reorder.Config{Fields: []string{"missing", "alpha"}})
	out := r.Apply(base())
	if len(out.Fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(out.Fields))
	}
	if out.Fields["alpha"] != "a" {
		t.Errorf("expected alpha=a, got %v", out.Fields["alpha"])
	}
}
