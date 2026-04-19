package clip_test

import (
	"testing"

	"logpipe/internal/clip"
)

func base() map[string]any {
	return map[string]any{
		"level":   "  info  ",
		"message": "\thello world\n",
		"count":   float64(3),
	}
}

func TestApply_NoConfig_TrimsAllStringFields(t *testing.T) {
	c := clip.New(clip.Config{})
	out := c.Apply(base())
	if got := out["level"]; got != "info" {
		t.Fatalf("level: got %q", got)
	}
	if got := out["message"]; got != "hello world" {
		t.Fatalf("message: got %q", got)
	}
	if got := out["count"]; got != float64(3) {
		t.Fatalf("count should be unchanged: got %v", got)
	}
}

func TestApply_SpecificField_OnlyTrimsTarget(t *testing.T) {
	c := clip.New(clip.Config{Fields: []string{"level"}})
	out := c.Apply(base())
	if got := out["level"]; got != "info" {
		t.Fatalf("level: got %q", got)
	}
	if got := out["message"]; got != "\thello world\n" {
		t.Fatalf("message should be untouched: got %q", got)
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	c := clip.New(clip.Config{Fields: []string{"LEVEL"}})
	out := c.Apply(base())
	if got := out["level"]; got != "info" {
		t.Fatalf("level: got %q", got)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	in := base()
	c := clip.New(clip.Config{})
	c.Apply(in)
	if got := in["level"]; got != "  info  " {
		t.Fatalf("original mutated: got %q", got)
	}
}

func TestApply_NonStringField_Unchanged(t *testing.T) {
	in := map[string]any{"count": float64(42)}
	c := clip.New(clip.Config{})
	out := c.Apply(in)
	if got := out["count"]; got != float64(42) {
		t.Fatalf("count: got %v", got)
	}
}
