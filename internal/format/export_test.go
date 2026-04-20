package format

// TargetOf exposes the configured target field for white-box tests.
func TargetOf(f *Formatter) string { return f.cfg.Target }

// TemplateOf exposes the configured template for white-box tests.
func TemplateOf(f *Formatter) string { return f.cfg.Template }
