package split_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/split"
)

type Entry = map[string]any

func base() Entry {
	return Entry{"level": "info", "msg": "hello", "tags": []any{"a", "b", "c"}}
}

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := split.New(split.Config{Field: ""})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_ArrayField_ExpandsEntries(t *testing.T) {
	s, _ := split.New(split.Config{Field: "tags"})
	out := s.Apply(base())
	if len(out) != 3 {
		t.Fatalf("want 3 entries, got %d", len(out))
	}
	for i, want := range []string{"a", "b", "c"} {
		if got := out[i]["tags"]; got != want {
			t.Errorf("entry %d: want tags=%q got %v", i, want, got)
		}
	}
}

func TestApply_CustomTargetField(t *testing.T) {
	s, _ := split.New(split.Config{Field: "tags", TargetField: "tag"})
	out := s.Apply(base())
	if len(out) != 3 {
		t.Fatalf("want 3 entries, got %d", len(out))
	}
	if _, ok := out[0]["tags"]; ok {
		t.Error("source field should be removed from child entries")
	}
	if out[0]["tag"] != "a" {
		t.Errorf("want tag=a got %v", out[0]["tag"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	s, _ := split.New(split.Config{Field: "TAGS"})
	out := s.Apply(base())
	if len(out) != 3 {
		t.Fatalf("want 3 entries, got %d", len(out))
	}
}

func TestApply_MissingField_PassesThrough(t *testing.T) {
	s, _ := split.New(split.Config{Field: "missing"})
	out := s.Apply(base())
	if len(out) != 1 {
		t.Fatalf("want 1 entry, got %d", len(out))
	}
}

func TestApply_NonSliceField_PassesThrough(t *testing.T) {
	s, _ := split.New(split.Config{Field: "msg"})
	out := s.Apply(base())
	if len(out) != 1 {
		t.Fatalf("want 1 entry, got %d", len(out))
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	s, _ := split.New(split.Config{Field: "tags"})
	e := base()
	s.Apply(e)
	if _, ok := e["tags"]; !ok {
		t.Error("original entry should not be mutated")
	}
}
