package sequence_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/sequence"
)

func base() map[string]any {
	return map[string]any{"message": "hello", "level": "info"}
}

func TestApply_DefaultField_InjectsSeq(t *testing.T) {
	s := sequence.New(sequence.Config{})
	out := s.Apply(base())
	if _, ok := out["_seq"]; !ok {
		t.Fatal("expected _seq field")
	}
}

func TestApply_MonotonicallyIncreases(t *testing.T) {
	s := sequence.New(sequence.Config{})
	a := s.Apply(base())
	b := s.Apply(base())
	va := a["_seq"].(int64)
	vb := b["_seq"].(int64)
	if vb != va+1 {
		t.Fatalf("expected %d+1 == %d", va, vb)
	}
}

func TestApply_CustomField_UsesField(t *testing.T) {
	s := sequence.New(sequence.Config{Field: "seq_num"})
	out := s.Apply(base())
	if _, ok := out["seq_num"]; !ok {
		t.Fatal("expected seq_num field")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	s := sequence.New(sequence.Config{})
	in := base()
	_ = s.Apply(in)
	if _, ok := in["_seq"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}

func TestReset_RestartsCounter(t *testing.T) {
	s := sequence.New(sequence.Config{})
	s.Apply(base())
	s.Apply(base())
	s.Reset()
	out := s.Apply(base())
	if out["_seq"].(int64) != 1 {
		t.Fatalf("expected counter to restart at 1, got %d", out["_seq"])
	}
}

func TestNew_DefaultField_IsSeq(t *testing.T) {
	s := sequence.New(sequence.Config{})
	if s.Field() != "_seq" {
		t.Fatalf("expected _seq, got %s", s.Field())
	}
}
