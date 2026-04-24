package fingerprint_test

import (
	"strings"
	"testing"

	"github.com/logpipe/internal/fingerprint"
)

func base() map[string]any {
	return map[string]any{
		"level":   "error",
		"message": "disk full",
		"host":    "web-01",
	}
}

func TestApply_DefaultOutputField(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{})
	out := fp.Apply(base())
	if _, ok := out["_fp"]; !ok {
		t.Fatal("expected _fp field to be present")
	}
}

func TestApply_CustomOutputField(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{OutputField: "fingerprint"})
	out := fp.Apply(base())
	if _, ok := out["fingerprint"]; !ok {
		t.Fatal("expected fingerprint field to be present")
	}
}

func TestApply_DeterministicForSameEntry(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{})
	a := fp.Apply(base())
	b := fp.Apply(base())
	if a["_fp"] != b["_fp"] {
		t.Fatalf("expected same fingerprint, got %v vs %v", a["_fp"], b["_fp"])
	}
}

func TestApply_DifferentEntries_DifferentFingerprints(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{})
	e1 := map[string]any{"message": "foo"}
	e2 := map[string]any{"message": "bar"}
	if fp.Apply(e1)["_fp"] == fp.Apply(e2)["_fp"] {
		t.Fatal("expected distinct fingerprints for distinct entries")
	}
}

func TestApply_SpecificFields_OnlyHashesThem(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{Fields: []string{"level"}})

	e1 := map[string]any{"level": "error", "message": "alpha"}
	e2 := map[string]any{"level": "error", "message": "beta"}

	if fp.Apply(e1)["_fp"] != fp.Apply(e2)["_fp"] {
		t.Fatal("fingerprints should match when only 'level' is used and both are 'error'")
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{
		Fields:          []string{"Level"},
		CaseInsensitive: true,
	})
	e := map[string]any{"level": "warn", "message": "x"}
	out := fp.Apply(e)
	if _, ok := out["_fp"]; !ok {
		t.Fatal("expected _fp to be set")
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{})
	orig := base()
	_ = fp.Apply(orig)
	if _, ok := orig["_fp"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}

func TestApply_FingerprintIsHexString(t *testing.T) {
	fp := fingerprint.New(fingerprint.Config{})
	out := fp.Apply(base())
	v, _ := out["_fp"].(string)
	if len(v) != 16 {
		t.Fatalf("expected 16-char hex fingerprint, got %q", v)
	}
	for _, c := range v {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Fatalf("non-hex character %q in fingerprint", c)
		}
	}
}
