package flatten_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/flatten"
)

func TestApply_NoNesting_PassesThrough(t *testing.T) {
	f := flatten.New(flatten.Config{})
	in := map[string]any{"level": "info", "msg": "hello"}
	out := f.Apply(in)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestApply_SingleNestedMap_IsFlattened(t *testing.T) {
	f := flatten.New(flatten.Config{})
	in := map[string]any{
		"level": "info",
		"http": map[string]any{"method": "GET", "status": 200},
	}
	out := f.Apply(in)
	if out["http.method"] != "GET" {
		t.Fatalf("expected http.method=GET, got %v", out["http.method"])
	}
	if out["http.status"] != 200 {
		t.Fatalf("expected http.status=200, got %v", out["http.status"])
	}
	if _, ok := out["http"]; ok {
		t.Fatal("parent key 'http' should be removed")
	}
}

func TestApply_DeepNesting_IsFlattened(t *testing.T) {
	f := flatten.New(flatten.Config{})
	in := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	out := f.Apply(in)
	if out["a.b.c"] != "deep" {
		t.Fatalf("expected a.b.c=deep, got %v", out)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := flatten.New(flatten.Config{Separator: "_"})
	in := map[string]any{
		"http": map[string]any{"method": "POST"},
	}
	out := f.Apply(in)
	if out["http_method"] != "POST" {
		t.Fatalf("expected http_method=POST, got %v", out)
	}
}

func TestApply_MaxDepth_LimitsRecursion(t *testing.T) {
	f := flatten.New(flatten.Config{MaxDepth: 1})
	inner := map[string]any{"c": "deep"}
	in := map[string]any{
		"a": map[string]any{
			"b": inner,
		},
	}
	out := f.Apply(in)
	// a.b should exist but b's children should not be further expanded
	if _, ok := out["a.b"]; !ok {
		t.Fatalf("expected a.b to be present, got %v", out)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	f := flatten.New(flatten.Config{})
	in := map[string]any{
		"meta": map[string]any{"env": "prod"},
	}
	_ = f.Apply(in)
	if _, ok := in["meta"]; !ok {
		t.Fatal("original entry was mutated")
	}
}
