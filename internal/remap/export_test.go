package remap

// RulesOf exposes internal rules for white-box tests.
func RulesOf(r *Remapper) []Rule {
	return r.rules
}
