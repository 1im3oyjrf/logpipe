package split

// FieldOf exposes the configured source field for white-box testing.
func FieldOf(s *Splitter) string { return s.field }

// TargetOf exposes the configured target field for white-box testing.
func TargetOf(s *Splitter) string { return s.target }
