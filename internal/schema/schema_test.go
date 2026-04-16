package schema_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/schema"
)

func TestValidate_RequiredFieldPresent(t *testing.T) {
	v := schema.New(schema.Config{
		RequiredFields: []string{"level", "message"},
		AllowUnknown:   true,
	})

	err := v.Validate(map[string]any{"level": "info", "message": "hello"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	v := schema.New(schema.Config{
		RequiredFields: []string{"level", "message"},
		AllowUnknown:   true,
	})

	err := v.Validate(map[string]any{"level": "info"})
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestValidate_CaseInsensitiveRequired(t *testing.T) {
	v := schema.New(schema.Config{
		RequiredFields: []string{"Level"},
		AllowUnknown:   true,
	})

	err := v.Validate(map[string]any{"level": "warn"})
	if err != nil {
		t.Fatalf("expected case-insensitive match, got %v", err)
	}
}

func TestValidate_UnknownFieldRejected(t *testing.T) {
	v := schema.New(schema.Config{
		AllowUnknown: false,
		KnownFields:  []string{"level", "message"},
	})

	err := v.Validate(map[string]any{"level": "info", "message": "ok", "extra": "bad"})
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestValidate_UnknownFieldAllowed(t *testing.T) {
	v := schema.New(schema.Config{
		AllowUnknown: true,
		KnownFields:  []string{"level", "message"},
	})

	err := v.Validate(map[string]any{"level": "info", "message": "ok", "extra": "fine"})
	if err != nil {
		t.Fatalf("expected no error when AllowUnknown=true, got %v", err)
	}
}

func TestValidate_EmptySchema_AcceptsAnything(t *testing.T) {
	v := schema.New(schema.Config{AllowUnknown: true})

	err := v.Validate(map[string]any{"foo": 1, "bar": "baz"})
	if err != nil {
		t.Fatalf("expected no error with empty schema, got %v", err)
	}
}
