package join_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/join"
)

func TestNew_DefaultTarget(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first"}})
	out := j.Apply(map[string]any{"first": "v"})
	if _, ok := out["joined"]; !ok {
		t.Fatal("expected default target field 'joined'")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "second"}})
	out := j.Apply(map[string]any{"first": "a", "second": "b"})
	if got := out["joined"]; got != "a b" {
		t.Fatalf("expected 'a b', got %q", got)
	}
}

func TestNew_ExplicitConfig_Respected(t *testing.T) {
	j := join.New(join.Config{
		Fields:    []string{"first", "second"},
		Separator: "|",
		Target:    "combined",
	})
	out := j.Apply(map[string]any{"first": "x", "second": "y"})
	if got := out["combined"]; got != "x|y" {
		t.Fatalf("expected 'x|y', got %q", got)
	}
}
