package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

const (
	pollInterval = 100 * time.Millisecond
)

// Tailer reads lines from a file, following new content as it is appended.
type Tailer struct {
	path   string
	lines  chan string
	errors chan error
}

// New creates a new Tailer for the given file path.
func New(path string) *Tailer {
	return &Tailer{
		path:   path,
		lines:  make(chan string, 64),
		errors: make(chan error, 1),
	}
}

// Lines returns the channel on which tailed lines are delivered.
func (t *Tailer) Lines() <-chan string {
	return t.lines
}

// Errors returns the channel on which errors are delivered.
func (t *Tailer) Errors() <-chan error {
	return t.errors
}

// Run starts tailing the file, sending lines to the Lines channel until ctx
// is cancelled or an unrecoverable error occurs.
func (t *Tailer) Run(ctx context.Context) {
	defer close(t.lines)
	defer close(t.errors)

	f, err := os.Open(t.path)
	if err != nil {
		t.errors <- err
		return
	}
	defer f.Close()

	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		t.errors <- err
		return
	}

	reader := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// No new data yet; wait before polling again.
				time.Sleep(pollInterval)
				continue
			}
			t.errors <- err
			return
		}

		if len(line) > 0 {
			select {
			case t.lines <- line:
			case <-ctx.Done():
				return
			}
		}
	}
}
