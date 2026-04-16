package truncate_test

import (
	"strings"
	"testing"

	"github.com/your-org/logpipe/internal/truncate"
)

func base() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": strings.Repeat("a", 300),
		"service": "svc",
	}
}

func TestApply_NoConfig_TruncatesAllLongFields(t *testing.T) {
	tr := truncate.New(truncate.Config{MaxLen: 10})
	out := tr.Apply(base())
	if got := out["message"].(string); len(got) != 10 {
		t.Fatalf("expected 10, got %d", len(got))
	}
	if got := out["level"].(string); got != "info" {
		t.Fatalf("short value mutated: %q", got)
	}
}

func TestApply_SpecificField_OnlyTruncatesThat(t *testing.T) {
	long := strings.Repeat("b", 50)
	entry := map[string]any{"message": long, "detail": long}
	tr := truncate.New(truncate.Config{MaxLen: 20, Fields: []string{"message"}})
	out := tr.Apply(entry)
	if got := out["message"].(string); len(got) != 20 {
		t.Fatalf("message: expected 20, got %d", len(got))
	}
	if got := out["detail"].(string); got != long {
		t.Fatalf("detail should be unchanged")
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	long := strings.Repeat("c", 50)
	entry := map[string]any{"Message": long}
	tr := truncate.New(truncate.Config{MaxLen: 5, Fields: []string{"message"}})
	out := tr.Apply(entry)
	if got := out["Message"].(string); len(got) != 5 {
		t.Fatalf("expected 5, got %d", len(got))
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	long := strings.Repeat("d", 100)
	entry := map[string]any{"msg": long}
	tr := truncate.New(truncate.Config{MaxLen: 10})
	_ = tr.Apply(entry)
	if got := entry["msg"].(string); len(got) != 100 {
		t.Fatal("original entry was mutated")
	}
}

func TestApply_DefaultMaxLen_UsedWhenZero(t *testing.T) {
	long := strings.Repeat("e", 300)
	entry := map[string]any{"msg": long}
	tr := truncate.New(truncate.Config{}) // MaxLen 0 → default 256
	out := tr.Apply(entry)
	if got := out["msg"].(string); len(got) != 256 {
		t.Fatalf("expected 256, got %d", len(got))
	}
}

func TestApply_NonStringField_Untouched(t *testing.T) {
	entry := map[string]any{"count": 42, "msg": strings.Repeat("f", 50)}
	tr := truncate.New(truncate.Config{MaxLen: 10})
	out := tr.Apply(entry)
	if got, ok := out["count"].(int); !ok || got != 42 {
		t.Fatal("numeric field should pass through unchanged")
	}
}
