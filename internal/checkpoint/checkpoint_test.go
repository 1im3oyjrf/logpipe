package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/logpipe/internal/checkpoint"
)

func TestNew_MissingFile_ReturnsEmptyStore(t *testing.T) {
	dir := t.TempDir()
	s, err := checkpoint.New(filepath.Join(dir, "cp.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("/var/log/app.log"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestSet_Get_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s, _ := checkpoint.New(filepath.Join(dir, "cp.json"))

	s.Set("/var/log/app.log", 1024)
	if got := s.Get("/var/log/app.log"); got != 1024 {
		t.Errorf("expected 1024, got %d", got)
	}
}

func TestFlush_PersistsToDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cp.json")

	s, _ := checkpoint.New(path)
	s.Set("/var/log/app.log", 2048)

	if err := s.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	s2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if got := s2.Get("/var/log/app.log"); got != 2048 {
		t.Errorf("expected 2048 after reload, got %d", got)
	}
}

func TestFlush_CreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "cp.json")

	s, _ := checkpoint.New(path)
	s.Set("src", 99)

	if err := s.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestFlush_OverwritesPreviousData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cp.json")

	s, _ := checkpoint.New(path)
	s.Set("a", 10)
	_ = s.Flush()

	s.Set("a", 20)
	s.Set("b", 30)
	_ = s.Flush()

	s2, _ := checkpoint.New(path)
	if got := s2.Get("a"); got != 20 {
		t.Errorf("expected 20, got %d", got)
	}
	if got := s2.Get("b"); got != 30 {
		t.Errorf("expected 30, got %d", got)
	}
}
