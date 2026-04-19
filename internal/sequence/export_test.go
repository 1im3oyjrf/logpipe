package sequence

// FieldOf exposes the configured field name for white-box tests.
func FieldOf(s *Sequencer) string { return s.field }

// CounterOf exposes the current counter value for white-box tests.
func CounterOf(s *Sequencer) uint64 { return s.counter }
