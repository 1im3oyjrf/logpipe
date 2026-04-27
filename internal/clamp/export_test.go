package clamp

// RulesOf exposes the internal rules slice for testing.
func RulesOf(c *Clamp) []Rule {
	return c.rules
}
