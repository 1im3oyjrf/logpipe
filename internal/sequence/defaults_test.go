package sequence_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/sequence"
)

func TestNew_ExplicitField_Respected(t *testing.T) {
	s := sequence.New(sequence.Config{Field: "n"})
	if s.Field() != "n" {
		t.Fatalf("expected n, got %s", s.Field())
	}
}

func TestNew_EmptyField_UsesDefault(t *testing.T) {
	s := sequence.New(sequence.Config{Field: ""})
	if s.Field() != "_seq" {
		t.Fatalf("expected _seq, got %s", s.Field())
	}
}

func TestNew_InitialCounter_IsZero(t *testing.T) {
	s := sequence.New(sequence.Config{})
	if sequence.CounterOf(s) != 0 {
		t.Fatal("initial counter must be zero")
	}
}
