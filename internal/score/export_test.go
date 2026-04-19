package score

// FieldOf exposes the configured source field for white-box tests.
func FieldOf(s *Scorer) string { return s.cfg.Field }

// OutputFieldOf exposes the configured output field for white-box tests.
func OutputFieldOf(s *Scorer) string { return s.cfg.OutputField }
