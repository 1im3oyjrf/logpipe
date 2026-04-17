package normalize_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/normalize"
)

func TestNew_DefaultLevelField(t *testing.T) {
	n := normalize.New(normalize.Config{})
	if normalize.LevelFieldOf(n) != "level" {
		t.Fatalf("expected default level field to be 'level'")
	}
}

func TestNew_FieldMap_KeysAreLowercased(t *testing.T) {
	n := normalize.New(normalize.Config{
		FieldMap: map[string]string{
			"Timestamp": "ts",
			"MESSAGE":   "msg",
		},
	})
	fm := normalize.FieldMapOf(n)
	if fm["timestamp"] != "ts" {
		t.Fatalf("expected 'timestamp' key, got map %v", fm)
	}
	if fm["message"] != "msg" {
		t.Fatalf("expected 'message' key, got map %v", fm)
	}
	if _, ok := fm["Timestamp"]; ok {
		t.Fatal("original mixed-case key should not be present")
	}
}
