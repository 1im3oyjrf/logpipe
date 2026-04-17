package timeout

// DroppedOf exposes the internal dropped counter for white-box tests.
func DroppedOf(g *Guard) uint64 { return g.dropped }

// DurationOf exposes the configured duration for white-box tests.
func DurationOf(g *Guard) interface{ String() string } { return g.cfg.Duration }
