package offset_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/offset"
)

func TestNew_DefaultTarget_UsesField(t *testing.T) {
	p := offset.New(offset.Config{Field: "score", By: 1})
	if offset.TargetOf(p) != "score" {
		t.Fatalf("expected target=score, got %s", offset.TargetOf(p))
	}
}

func TestNew_ExplicitTarget_Stored(t *testing.T) {
	p := offset.New(offset.Config{Field: "score", By: 1, Target: "score_adj"})
	if offset.TargetOf(p) != "score_adj" {
		t.Fatalf("expected target=score_adj, got %s", offset.TargetOf(p))
	}
}

func TestNew_ByStored(t *testing.T) {
	p := offset.New(offset.Config{Field: "x", By: -7.5})
	if offset.ByOf(p) != -7.5 {
		t.Fatalf("expected by=-7.5, got %v", offset.ByOf(p))
	}
}
