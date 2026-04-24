// Package expire provides an Expirer that drops log entries whose timestamp
// field indicates the entry is older than a configured maximum age.
//
// # Usage
//
//	cfg := expire.Config{
//		TimestampField: "ts",   // field holding RFC3339 timestamp (default: "timestamp")
//		MaxAge:         time.Hour, // entries older than this are dropped (default: 5m)
//	}
//	e, err := expire.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if e.Allow(entry) {
//		// forward entry downstream
//	}
//
// Entries with a missing or unparseable timestamp are always forwarded so that
// non-timestamped log lines are not silently discarded.
package expire
