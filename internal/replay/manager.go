package replay

import (
	"context"
	"sync"

	"github.com/logpipe/internal/buffer"
	"github.com/logpipe/internal/reader"
)

// Manager records incoming entries into a ring buffer and
// provides on-demand replay of historical entries.
type Manager struct {
	mu  sync.Mutex
	buf *buffer.Ring
}

// NewManager creates a Manager with the given buffer capacity.
func NewManager(capacity int) *Manager {
	return &Manager{buf: buffer.New(capacity)}
}

// Record stores an entry in the ring buffer.
func (m *Manager) Record(e reader.Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.buf.Push(e)
}

// Replay returns a channel of historical entries matching cfg.
func (m *Manager) Replay(ctx context.Context, cfg Config) <-chan reader.Entry {
	m.mu.Lock()
	r := New(m.buf, cfg)
	m.mu.Unlock()
	return r.Replay(ctx)
}

// Size returns the number of entries currently held in the buffer.
func (m *Manager) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.buf.Snapshot())
}

// Reset clears all entries from the ring buffer, discarding recorded history.
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.buf.Clear()
}
