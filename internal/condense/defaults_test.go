package condense_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/condense"
)

func TestNew_DefaultCountField(t *testing.T) {
	c := condense.New(condense.Config{})
	if condense.CountFieldOf(c) != "_repeat" {
		t.Fatalf("expected default count field _repeat, got %q", condense.CountFieldOf(c))
	}
}

func TestNew_ExplicitCountField_Respected(t *testing.T) {
	c := condense.New(condense.Config{CountField: "n"})
	if condense.CountFieldOf(c) != "n" {
		t.Fatalf("expected count field n, got %q", condense.CountFieldOf(c))
	}
}
