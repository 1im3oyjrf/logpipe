package join_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/join"
)

var base = map[string]any{
	"first": "hello",
	"second": "world",
	"level":  "info",
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	j := join.New(join.Config{})
	out := j.Apply(base)
	if _, ok := out["joined"]; ok {
		t.Fatal("expected no joined field when Fields is empty")
	}
}

func TestApply_JoinsTwoFields(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "second"}})
	out := j.Apply(base)
	if got := out["joined"]; got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "second"}, Separator: "-"})
	out := j.Apply(base)
	if got := out["joined"]; got != "hello-world" {
		t.Fatalf("expected 'hello-world', got %q", got)
	}
}

func TestApply_CustomTarget(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "second"}, Target: "msg"})
	out := j.Apply(base)
	if _, ok := out["msg"]; !ok {
		t.Fatal("expected field 'msg' in output")
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "missing"}})
	out := j.Apply(base)
	if got := out["joined"]; got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestApply_DropSources_RemovesFields(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "second"}, DropSources: true})
	out := j.Apply(base)
	if _, ok := out["first"]; ok {
		t.Fatal("expected 'first' to be removed")
	}
	if _, ok := out["second"]; ok {
		t.Fatal("expected 'second' to be removed")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	j := join.New(join.Config{Fields: []string{"first", "second"}, DropSources: true})
	_ = j.Apply(base)
	if _, ok := base["first"]; !ok {
		t.Fatal("original entry must not be mutated")
	}
}
