package join

// SeparatorOf exposes the configured separator for white-box tests.
func SeparatorOf(j *Joiner) string { return j.cfg.Separator }

// TargetOf exposes the configured target field for white-box tests.
func TargetOf(j *Joiner) string { return j.cfg.Target }
