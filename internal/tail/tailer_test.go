package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/logpipe/internal/tail"
)

func writeTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return f
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	f := writeTempFile(t, "")
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tr := tail.New(f.Name())
	go tr.Run(ctx)

	// Give the tailer time to reach EOF before writing.
	time.Sleep(50 * time.Millisecond)

	want := `{"level":"info","msg":"hello"}` + "\n"
	if _, err := f.WriteString(want); err != nil {
		t.Fatalf("append to file: %v", err)
	}

	select {
	case got := <-tr.Lines():
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	case err := <-tr.Errors():
		t.Fatalf("unexpected error: %v", err)
	case <-ctx.Done():
		t.Fatal("timed out waiting for line")
	}
}

func TestTailer_MissingFile_ReportsError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tr := tail.New("/nonexistent/path/to/file.log")
	go tr.Run(ctx)

	select {
	case err := <-tr.Errors():
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for error")
	}
}

func TestTailer_ContextCancellation_StopsGracefully(t *testing.T) {
	f := writeTempFile(t, "")
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())

	tr := tail.New(f.Name())
	done := make(chan struct{})
	go func() {
		tr.Run(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Fatal("tailer did not stop after context cancellation")
	}
}
