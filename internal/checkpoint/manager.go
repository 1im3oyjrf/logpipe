package checkpoint

import (
	"context"
	"time"
)

// Manager periodically flushes a Store to disk and provides a clean
// shutdown path via its Run method.
type Manager struct {
	store    *Store
	interval time.Duration
}

// NewManager creates a Manager that will flush store every interval.
// A zero or negative interval defaults to 10 seconds.
func NewManager(store *Store, interval time.Duration) *Manager {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	return &Manager{store: store, interval: interval}
}

// Run blocks, flushing the store on every tick until ctx is cancelled,
// at which point it performs one final flush before returning.
func (m *Manager) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.store.Flush(); err != nil {
				return err
			}
		case <-ctx.Done():
			return m.store.Flush()
		}
	}
}
