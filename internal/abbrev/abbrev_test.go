package abbrev_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/abbrev"
	"github.com/yourorg/logpipe/internal/reader"
)

var base = reader.Entry{
	"level":   "info",
	"message": "this is a very long message that should definitely be truncated by the abbreviator",
	"short":   "ok",
}

func TestApply_NoConfig_TruncatesAllLongFields(t *testing.T) {
	a := abbrev.New(abbrev.Config{})
	out := a.Apply(base)
	msg, _ := out["message"].(string)
	if !strings.HasSuffix(msg, "...") {
		t.Fatalf("expected truncated suffix, got %q", msg)
	}
	if len([]rune(msg)) != 80+3 {
		t.Fatalf("expected length %d, got %d", 80+3, len([]rune(msg)))
	}
}

func TestApply_ShortField_NotTruncated(t *testing.T) {
	a := abbrev.New(abbrev.Config{})
	out := a.Apply(base)
	if out["short"] != "ok" {
		t.Fatalf("short field should not be modified, got %v", out["short"])
	}
}

func TestApply_SpecificField_OnlyAbbreviatesThat(t *testing.T) {
	a := abbrev.New(abbrev.Config{Fields: []string{"message"}, MaxLen: 10})
	out := a.Apply(base)
	level, _ := out["level"].(string)
	if level != "info" {
		t.Fatalf("level should be untouched, got %q", level)
	}
	msg, _ := out["message"].(string)
	if len([]rune(msg)) != 10+3 {
		t.Fatalf("expected abbreviated length %d, got %d", 13, len([]rune(msg)))
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	a := abbrev.New(abbrev.Config{
		Fields:          []string{"MESSAGE"},
		MaxLen:          10,
		CaseInsensitive: true,
	})
	out := a.Apply(base)
	msg, _ := out["message"].(string)
	if !strings.HasSuffix(msg, "...") {
		t.Fatalf("expected truncation on case-insensitive match, got %q", msg)
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	a := abbrev.New(abbrev.Config{MaxLen: 5, Suffix: "[…]"})
	entry := reader.Entry{"msg": "hello world"}
	out := a.Apply(entry)
	v, _ := out["msg"].(string)
	if !strings.HasSuffix(v, "[…]") {
		t.Fatalf("expected custom suffix, got %q", v)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	a := abbrev.New(abbrev.Config{MaxLen: 5})
	original := "hello world"
	entry := reader.Entry{"msg": original}
	a.Apply(entry)
	if entry["msg"] != original {
		t.Fatal("original entry should not be mutated")
	}
}

func TestApply_NonStringField_Skipped(t *testing.T) {
	a := abbrev.New(abbrev.Config{MaxLen: 1})
	entry := reader.Entry{"count": 42}
	out := a.Apply(entry)
	if out["count"] != 42 {
		t.Fatalf("non-string field should be unchanged, got %v", out["count"])
	}
}
