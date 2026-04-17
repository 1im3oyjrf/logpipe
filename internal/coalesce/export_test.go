package coalesce

// RulesOf exposes the internal rules slice for white-box tests.
func RulesOf(t *Transformer) []Rule {
	return t.rules
}
