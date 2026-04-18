package pivot_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/pivot"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
		"key":     "region",
		"value":   "us-east-1",
	}
}

func TestApply_NoKeyField_PassesThrough(t *testing.T) {
	p := pivot.New(pivot.Config{KeyField: "missing", ValueField: "value"})
	in := base()
	out := p.Apply(in)
	if _, ok := out["region"]; ok {
		t.Fatal("expected no pivot when key field absent")
	}
}

func TestApply_PivotCreatesNewField(t *testing.T) {
	p := pivot.New(pivot.Config{KeyField: "key", ValueField: "value"})
	out := p.Apply(base())
	if out["region"] != "us-east-1" {
		t.Fatalf("expected region=us-east-1, got %v", out["region"])
	}
}

func TestApply_DropSource_RemovesKeyAndValue(t *testing.T) {
	p := pivot.New(pivot.Config{KeyField: "key", ValueField: "value", DropSource: true})
	out := p.Apply(base())
	if _, ok := out["key"]; ok {
		t.Fatal("key field should be removed")
	}
	if _, ok := out["value"]; ok {
		t.Fatal("value field should be removed")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	p := pivot.New(pivot.Config{KeyField: "key", ValueField: "value", DropSource: true})
	in := base()
	p.Apply(in)
	if _, ok := in["key"]; !ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestApply_CaseInsensitiveKeyField(t *testing.T) {
	p := pivot.New(pivot.Config{KeyField: "KEY", ValueField: "VALUE"})
	out := p.Apply(base())
	if out["region"] != "us-east-1" {
		t.Fatalf("case-insensitive pivot failed, got %v", out["region"])
	}
}

func TestApply_DefaultFieldNames(t *testing.T) {
	p := pivot.New(pivot.Config{})
	out := p.Apply(base())
	if out["region"] != "us-east-1" {
		t.Fatalf("default field names failed, got %v", out["region"])
	}
}
