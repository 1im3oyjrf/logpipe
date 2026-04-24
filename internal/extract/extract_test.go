package extract_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/extract"
)

func base() map[string]any {
	return map[string]any{
		"level": "info",
		"message": "hello",
		"metadata": map[string]any{
			"request_id": "abc-123",
			"user_id":    float64(42),
		},
	}
}

func TestNew_EmptyPaths_ReturnsError(t *testing.T) {
	_, err := extract.New(extract.Config{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestApply_NoMatchingPath_PassesThrough(t *testing.T) {
	ex, _ := extract.New(extract.Config{Paths: []string{"metadata.missing"}})
	out := ex.Apply(base())
	if _, ok := out["metadata.missing"]; ok {
		t.Fatal("unexpected key in output")
	}
	if out["level"] != "info" {
		t.Fatal("original fields should be preserved")
	}
}

func TestApply_ExtractsNestedField(t *testing.T) {
	ex, _ := extract.New(extract.Config{Paths: []string{"metadata.request_id"}})
	out := ex.Apply(base())
	if got, ok := out["metadata.request_id"]; !ok || got != "abc-123" {
		t.Fatalf("expected 'abc-123', got %v", got)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	in := base()
	ex, _ := extract.New(extract.Config{Paths: []string{"metadata.request_id"}})
	ex.Apply(in)
	if _, ok := in["metadata.request_id"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestApply_DropSource_RemovesLeaf(t *testing.T) {
	ex, _ := extract.New(extract.Config{
		Paths:      []string{"metadata.request_id"},
		DropSource: true,
	})
	out := ex.Apply(base())
	nested, ok := out["metadata"].(map[string]any)
	if !ok {
		t.Fatal("metadata should still be present")
	}
	if _, found := nested["request_id"]; found {
		t.Fatal("request_id should have been removed from nested map")
	}
	if nested["user_id"] != float64(42) {
		t.Fatal("other nested fields should be preserved")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	ex, _ := extract.New(extract.Config{
		Paths:           []string{"Metadata.Request_ID"},
		CaseInsensitive: true,
	})
	out := ex.Apply(base())
	if got, ok := out["Metadata.Request_ID"]; !ok || got != "abc-123" {
		t.Fatalf("expected 'abc-123' via case-insensitive match, got %v", got)
	}
}

func TestApply_MultiplePaths_AllExtracted(t *testing.T) {
	ex, _ := extract.New(extract.Config{
		Paths: []string{"metadata.request_id", "metadata.user_id"},
	})
	out := ex.Apply(base())
	if out["metadata.request_id"] != "abc-123" {
		t.Fatal("request_id not extracted")
	}
	if out["metadata.user_id"] != float64(42) {
		t.Fatal("user_id not extracted")
	}
}
