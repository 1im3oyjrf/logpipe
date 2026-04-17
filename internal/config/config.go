package config

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds the parsed CLI configuration for logpipe.
type Config struct {
	// Sources is a list of file paths or "-" for stdin.
	Sources []string

	// Pattern is the optional grep/filter pattern.
	Pattern string

	// CaseSensitive controls whether pattern matching is case-sensitive.
	CaseSensitive bool

	// Fields restricts output to a specific set of JSON fields.
	Fields []string

	// NoColor disables ANSI color output.
	NoColor bool

	// Level filters log entries to only those at or above the given level.
	Level string
}

// validLevels contains the accepted log level values.
var validLevels = map[string]bool{
	"":      true,
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

// Parse reads CLI flags and arguments, returning a populated Config.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logpipe", flag.ContinueOnError)

	pattern := fs.String("grep", "", "filter pattern (substring or regex)")
	caseSensitive := fs.Bool("case-sensitive", false, "enable case-sensitive matching")
	fieldsRaw := fs.String("fields", "", "comma-separated list of fields to display")
	noColor := fs.Bool("no-color", false, "disable color output")
	level := fs.String("level", "", "minimum log level to display (debug, info, warn, error)")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	normalizedLevel := strings.ToLower(*level)
	if !validLevels[normalizedLevel] {
		return nil, fmt.Errorf("invalid level %q: must be one of debug, info, warn, error", *level)
	}

	sources := fs.Args()
	if len(sources) == 0 {
		sources = []string{"-"}
	}

	var fields []string
	if *fieldsRaw != "" {
		for _, f := range strings.Split(*fieldsRaw, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				fields = append(fields, f)
			}
		}
	}

	return &Config{
		Sources:       sources,
		Pattern:       *pattern,
		CaseSensitive: *caseSensitive,
		Fields:        fields,
		NoColor:       *noColor,
		Level:         normalizedLevel,
	}, nil
}
