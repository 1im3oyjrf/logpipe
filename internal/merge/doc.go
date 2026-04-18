// Package merge provides the Merger processor which injects a static set of
// key/value fields into every log entry passing through the pipeline.
//
// Fields are matched case-insensitively; the stored keys are always
// lower-cased so downstream processors see a consistent schema.
//
// Example usage:
//
//	m := merge.New(merge.Config{
//		Fields:    map[string]string{"env": "production", "region": "us-east-1"},
//		Overwrite: false,
//	})
//	enriched := m.Apply(entry)
package merge
