package count_test

import (
	"strconv"
	"testing"

	"github.com/your-org/logpipe/internal/count"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello",
	}
}

func TestApply_DefaultField_InjectsCount(t *testing.T) {
	p := count.New(count.Config{})
	out := p.Apply(base())
	if _, ok := out["_count"]; !ok {
		t.Fatal("expected _count field to be present")
	}
}

func TestApply_CounterIncrementsEachCall(t *testing.T) {
	p := count.New(count.Config{})
	for i := 1; i <= 5; i++ {
		out := p.Apply(base())
		got, _ := strconv.Atoi(out["_count"].(string))
		if got != i {
			t.Fatalf("call %d: expected counter %d, got %d", i, i, got)
		}
	}
}

func TestApply_CustomField_UsesField(t *testing.T) {
	p := count.New(count.Config{Field: "seq"})
	out := p.Apply(base())
	if _, ok := out["seq"]; !ok {
		t.Fatal("expected seq field to be present")
	}
	if _, ok := out["_count"]; ok {
		t.Fatal("unexpected _count field")
	}
}

func TestApply_ExistingField_NotOverwrittenByDefault(t *testing.T) {
	p := count.New(count.Config{})
	entry := base()
	entry["_count"] = "99"
	out := p.Apply(entry)
	if out["_count"] != "99" {
		t.Fatalf("expected original value to be preserved, got %v", out["_count"])
	}
}

func TestApply_Overwrite_ReplacesExistingField(t *testing.T) {
	p := count.New(count.Config{Overwrite: true})
	entry := base()
	entry["_count"] = "99"
	out := p.Apply(entry)
	if out["_count"] == "99" {
		t.Fatal("expected field to be overwritten")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	p := count.New(count.Config{})
	original := base()
	p.Apply(original)
	if _, ok := original["_count"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	p := count.New(count.Config{})
	p.Apply(base())
	p.Apply(base())
	if p.Value() != 2 {
		t.Fatalf("expected value 2 before reset, got %d", p.Value())
	}
	p.Reset()
	if p.Value() != 0 {
		t.Fatalf("expected value 0 after reset, got %d", p.Value())
	}
	out := p.Apply(base())
	if out["_count"] != "1" {
		t.Fatalf("expected counter to restart at 1, got %v", out["_count"])
	}
}
