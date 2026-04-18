package ceiling

// WindowOf exposes the configured window duration for testing.
func WindowOf(c *Ceiling) time.Duration {
	return c.cfg.Window
}

// MaxOf exposes the configured max for testing.
func MaxOf(c *Ceiling) int {
	return c.cfg.Max
}
