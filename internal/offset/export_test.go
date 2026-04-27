package offset

// TargetOf exposes the resolved target field for white-box tests.
func TargetOf(p *Processor) string {
	return p.cfg.Target
}

// ByOf exposes the shift amount for white-box tests.
func ByOf(p *Processor) float64 {
	return p.cfg.By
}
