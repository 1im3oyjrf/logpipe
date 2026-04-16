package multiwriter_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/logpipe/logpipe/internal/multiwriter"
)

type errWriter struct{ err error }

func (e *errWriter) Write(_ []byte) (int, error) { return 0, e.err }

func TestNew_WritesToAllTargets(t *testing.T) {
	var a, b bytes.Buffer
	w := multiwriter.New(&a, &b)
	w.Write([]byte("hello"))
	if a.String() != "hello" || b.String() != "hello" {
		t.Fatalf("expected both buffers to contain 'hello', got %q %q", a.String(), b.String())
	}
}

func TestAdd_AppendsTarget(t *testing.T) {
	var a, b bytes.Buffer
	w := multiwriter.New(&a)
	w.Add(&b)
	if w.Len() != 2 {
		t.Fatalf("expected 2 targets, got %d", w.Len())
	}
	w.Write([]byte("x"))
	if b.String() != "x" {
		t.Fatalf("expected 'x' in added target, got %q", b.String())
	}
}

func TestRemove_DeletesTarget(t *testing.T) {
	var a, b bytes.Buffer
	w := multiwriter.New(&a, &b)
	w.Remove(&a)
	if w.Len() != 1 {
		t.Fatalf("expected 1 target after remove, got %d", w.Len())
	}
	w.Write([]byte("y"))
	if a.String() != "" {
		t.Fatalf("removed target should not receive writes, got %q", a.String())
	}
	if b.String() != "y" {
		t.Fatalf("remaining target should receive write, got %q", b.String())
	}
}

func TestWrite_ReturnsErrorFromFailingTarget(t *testing.T) {
	var buf bytes.Buffer
	ew := &errWriter{err: errors.New("write failed")}
	w := multiwriter.New(&buf, ew)
	_, err := w.Write([]byte("data"))
	if err == nil {
		t.Fatal("expected error from failing writer")
	}
	if buf.String() != "data" {
		t.Fatalf("healthy target should still receive write, got %q", buf.String())
	}
}

func TestWrite_EmptyTargets_NoError(t *testing.T) {
	w := multiwriter.New()
	_, err := w.Write([]byte("noop"))
	if err != nil {
		t.Fatalf("unexpected error with no targets: %v", err)
	}
}

func TestWrite_ImplementsIOWriter(t *testing.T) {
	w := multiwriter.New()
	var _ io.Writer = w
}
