// Package clamp provides the Clamp processor which constrains the values of
// numeric fields in a log entry to a configured [Min, Max] range.
//
// Values below Min are raised to Min; values above Max are lowered to Max.
// Fields that are absent or non-numeric are left untouched.
//
// Example usage:
//
//	c, err := clamp.New(clamp.Config{
//		Rules: []clamp.Rule{
//			{Field: "duration_ms", Min: 0, Max: 30000},
//			{Field: "score",       Min: 1, Max: 10},
//		},
//		CaseInsensitive: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	processed := c.Apply(entry)
package clamp
