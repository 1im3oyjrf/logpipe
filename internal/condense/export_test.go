package condense

// CountFieldOf exposes the configured count field name for white-box tests.
func CountFieldOf(c *Condenser) string {
	return c.cfg.CountField
}
