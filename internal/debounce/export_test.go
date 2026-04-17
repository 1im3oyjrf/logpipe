package debounce

// QuietPeriodOf exposes the resolved quiet period for white-box tests.
func QuietPeriodOf(d *Debouncer) interface{} {
	return d.cfg.QuietPeriod
}

// KeyFieldOf exposes the resolved key field for white-box tests.
func KeyFieldOf(d *Debouncer) string {
	return d.cfg.KeyField
}
