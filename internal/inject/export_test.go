package inject

// FieldsOf exposes the normalised fields map for white-box testing.
func FieldsOf(inj *Injector) map[string]string {
	return inj.cfg.Fields
}
