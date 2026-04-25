package fold

// SeparatorOf exposes the configured separator for white-box testing.
func SeparatorOf(f *Folder) string {
	return f.cfg.Separator
}

// TargetOf exposes the configured target field for white-box testing.
func TargetOf(f *Folder) string {
	return f.cfg.Target
}
