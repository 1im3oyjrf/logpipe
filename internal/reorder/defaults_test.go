package reorder_test

import (
	"testing"

	"logpipe/internal/reorder"
)

func TestNew_EmptyConfig_NoFields(t *testing.T) {
	r := reorder.New(reorder.Config{})
	if fields := reorder.FieldsOf(r); len(fields) != 0 {
		t.Errorf("expected no fields, got %v", fields)
	}
}

func TestNew_ExplicitFields_Stored(t *testing.T) {
	cfg := reorder.Config{Fields: []string{"msg", "level"}}
	r := reorder.New(cfg)
	got := reorder.FieldsOf(r)
	if len(got) != 2 || got[0] != "msg" || got[1] != "level" {
		t.Errorf("unexpected fields: %v", got)
	}
}

func TestNew_CaseInsensitiveDefault_IsFalse(t *testing.T) {
	r := reorder.New(reorder.Config{})
	if reorder.CaseInsensitiveOf(r) {
		t.Error("expected CaseInsensitive to default to false")
	}
}

func TestNew_ExplicitCaseInsensitive_Stored(t *testing.T) {
	r := reorder.New(reorder.Config{CaseInsensitive: true})
	if !reorder.CaseInsensitiveOf(r) {
		t.Error("expected CaseInsensitive to be true")
	}
}
