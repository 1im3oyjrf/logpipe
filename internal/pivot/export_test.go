package pivot

// KeyFieldOf exposes the configured key field for white-box tests.
func KeyFieldOf(p *Pivoter) string { return p.cfg.KeyField }

// ValueFieldOf exposes the configured value field for white-box tests.
func ValueFieldOf(p *Pivoter) string { return p.cfg.ValueField }
