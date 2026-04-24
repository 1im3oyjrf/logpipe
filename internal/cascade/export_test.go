package cascade

// RulesOf exposes the internal rules slice for white-box testing.
func RulesOf(c *Cascade) []Rule {
	return c.rules
}

// StopOnFirstOf exposes the stopOnFirst flag for white-box testing.
func StopOnFirstOf(c *Cascade) bool {
	return c.stopOnFirst
}
