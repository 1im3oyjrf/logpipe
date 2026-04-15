// export_test.go exposes internal fields for white-box testing.
package throttle

// SetClock replaces the internal clock function used by Throttle.
// This is only available during testing.
func (t *Throttle) SetClock(fn func() time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.now = fn
}
