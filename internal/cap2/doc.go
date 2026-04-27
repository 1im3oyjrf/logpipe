// Package cap2 provides a per-value entry cap for log pipelines.
//
// Unlike the cap package which applies a global limit across all entries,
// cap2 tracks counts independently per distinct field value. This allows
// capping the number of log entries emitted for each unique value of a
// given field — for example, limiting to 5 errors per unique error code.
//
// Usage:
//
//	 c, err := cap2.New(cap2.Config{
//	     Field: "error_code",
//	     Max:   5,
//	 })
//	 if err != nil {
//	     log.Fatal(err)
//	 }
//	 if c.Allow(entry) {
//	     // forward entry
//	 }
package cap2
