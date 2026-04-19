package reorder

// FieldsOf exposes the configured field list for white-box tests.
func FieldsOf(r *Reorder) []string {
	return r.cfg.Fields
}

// CaseInsensitiveOf exposes the case-insensitive flag for white-box tests.
func CaseInsensitiveOf(r *Reorder) bool {
	return r.cfg.CaseInsensitive
}
