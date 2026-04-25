package fold

import "testing"

func TestNew_DefaultSeparator_IsSpace(t *testing.T) {
	f, err := New(Config{Fields: []string{"a", "b"}, Target: "out"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := SeparatorOf(f); got != " " {
		t.Errorf("want " "", got %q", got)
	}
}

func TestNew_ExplicitSeparator_Respected(t *testing.T) {
	f, err := New(Config{Fields: []string{"a", "b"}, Target: "out", Separator: "|"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := SeparatorOf(f); got != "|" {
		t.Errorf("want "|"", got %q", got)
	}
}

func TestNew_TargetStored(t *testing.T) {
	f, err := New(Config{Fields: []string{"a", "b"}, Target: "combined"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := TargetOf(f); got != "combined" {
		t.Errorf("want %q, got %q", "combined", got)
	}
}
