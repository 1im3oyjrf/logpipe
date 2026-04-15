package config_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func TestParse_Defaults(t *testing.T) {
	cfg, err := config.Parse([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 1 || cfg.Sources[0] != "-" {
		t.Errorf("expected default source stdin, got %v", cfg.Sources)
	}
	if cfg.Pattern != "" {
		t.Errorf("expected empty pattern, got %q", cfg.Pattern)
	}
	if cfg.NoColor {
		t.Error("expected NoColor to be false by default")
	}
}

func TestParse_WithPattern(t *testing.T) {
	cfg, err := config.Parse([]string{"-grep", "ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pattern != "ERROR" {
		t.Errorf("expected pattern ERROR, got %q", cfg.Pattern)
	}
}

func TestParse_WithFields(t *testing.T) {
	cfg, err := config.Parse([]string{"-fields", "level, msg, ts"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(cfg.Fields))
	}
	if cfg.Fields[1] != "msg" {
		t.Errorf("expected second field 'msg', got %q", cfg.Fields[1])
	}
}

func TestParse_MultipleSourceFiles(t *testing.T) {
	cfg, err := config.Parse([]string{"app.log", "worker.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(cfg.Sources))
	}
}

func TestParse_LevelNormalized(t *testing.T) {
	cfg, err := config.Parse([]string{"-level", "WARN"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Level != "warn" {
		t.Errorf("expected level 'warn', got %q", cfg.Level)
	}
}

func TestParse_NoColorFlag(t *testing.T) {
	cfg, err := config.Parse([]string{"-no-color"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("expected NoColor to be true")
	}
}
