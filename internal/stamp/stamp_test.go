package stamp_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/stamp"
)

func base() map[string]any {
	return map[string]any{"message": "hello", "level": "info"}
}

func TestApply_NoConfig_InjectsTimestamp(t *testing.T) {
	s := stamp.New(stamp.Config{})
	out := s.Apply(base())
	if _, ok := out["timestamp"]; !ok {
		t.Fatal("expected timestamp field to be injected")
	}
}

func TestApply_CustomField_UsesField(t *testing.T) {
	s := stamp.New(stamp.Config{Field: "ts"})
	out := s.Apply(base())
	if _, ok := out["ts"]; !ok {
		t.Fatal("expected ts field")
	}
}

func TestApply_ExistingField_NotOverwritten(t *testing.T) {
	s := stamp.New(stamp.Config{Overwrite: false})
	in := base()
	in["timestamp"] = "original"
	out := s.Apply(in)
	if out["timestamp"] != "original" {
		t.Fatalf("expected original, got %v", out["timestamp"])
	}
}

func TestApply_Overwrite_ReplacesField(t *testing.T) {
	s := stamp.New(stamp.Config{Overwrite: true})
	in := base()
	in["timestamp"] = "old"
	out := s.Apply(in)
	if out["timestamp"] == "old" {
		t.Fatal("expected timestamp to be overwritten")
	}
}

func TestApply_CustomFormat(t *testing.T) {
	const layout = "2006-01-02"
	s := stamp.New(stamp.Config{Format: layout})
	out := s.Apply(base())
	v, _ := out["timestamp"].(string)
	if _, err := time.Parse(layout, v); err != nil {
		t.Fatalf("timestamp %q does not match layout: %v", v, err)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	s := stamp.New(stamp.Config{})
	in := base()
	_ = s.Apply(in)
	if _, ok := in["timestamp"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}
