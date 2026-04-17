package normalize

// LevelFieldOf exposes the resolved level field name for white-box tests.
func LevelFieldOf(n *Normalizer) string {
	return n.cfg.LevelField
}

// FieldMapOf exposes the normalised field map for white-box tests.
func FieldMapOf(n *Normalizer) map[string]string {
	copy := make(map[string]string, len(n.cfg.FieldMap))
	for k, v := range n.cfg.FieldMap {
		copy[k] = v
	}
	return copy
}
