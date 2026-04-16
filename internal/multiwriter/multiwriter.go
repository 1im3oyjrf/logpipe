package multiwriter

import (
	"io"
	"sync"
)

// Writer fans out writes to multiple underlying io.Writer targets.
type Writer struct {
	mu      sync.RWMutex
	targets []io.Writer
}

// New returns a Writer that broadcasts to the given targets.
func New(targets ...io.Writer) *Writer {
	out := make([]io.Writer, len(targets))
	copy(out, targets)
	return &Writer{targets: out}
}

// Add appends a new target to the writer.
func (w *Writer) Add(t io.Writer) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.targets = append(w.targets, t)
}

// Remove removes the first target that matches t by pointer identity.
func (w *Writer) Remove(t io.Writer) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i, target := range w.targets {
		if target == t {
			w.targets = append(w.targets[:i], w.targets[i+1:]...)
			return
		}
	}
}

// Write writes p to all targets. It collects all errors and returns the
// last non-nil error encountered, along with the byte count from the
// first successful write (or 0 on total failure).
func (w *Writer) Write(p []byte) (int, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var (
		n   int
		err error
	)
	for _, t := range w.targets {
		written, werr := t.Write(p)
		if werr != nil {
			err = werr
		} else if n == 0 {
			n = written
		}
	}
	return n, err
}

// Len returns the current number of registered targets.
func (w *Writer) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.targets)
}
