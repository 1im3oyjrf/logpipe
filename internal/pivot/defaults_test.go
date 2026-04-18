package pivot_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/pivot"
)

func TestNew_DefaultKeyField(t *testing.T) {
	p := pivot.New(pivot.Config{})
	if pivot.KeyFieldOf(p) != "key" {
		t.Fatalf("expected default key field 'key', got %q", pivot.KeyFieldOf(p))
	}
}

func TestNew_DefaultValueField(t *testing.T) {
	p := pivot.New(pivot.Config{})
	if pivot.ValueFieldOf(p) != "value" {
		t.Fatalf("expected default value field 'value', got %q", pivot.ValueFieldOf(p))
	}
}

func TestNew_ExplicitFields_Respected(t *testing.T) {
	p := pivot.New(pivot.Config{KeyField: "attr", ValueField: "data"})
	if pivot.KeyFieldOf(p) != "attr" {
		t.Fatalf("expected key field 'attr', got %q", pivot.KeyFieldOf(p))
	}
	if pivot.ValueFieldOf(p) != "data" {
		t.Fatalf("expected value field 'data', got %q", pivot.ValueFieldOf(p))
	}
}
