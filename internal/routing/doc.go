// Package routing implements rule-based dispatch of log entries to named
// output channels. Each Rule specifies a field name, an expected value, and
// a target channel. The Router evaluates rules in order and sends the entry
// to the first matching target. Entries that satisfy no rule are forwarded
// to the built-in "default" channel.
//
// Example usage:
//
//	rules := []routing.Rule{
//		{Field: "level", Value: "error", Target: "errors"},
//		{Field: "service", Value: "auth",  Target: "auth"},
//	}
//	router := routing.New(rules, 64)
//	defer router.Close()
//
//	go consumeChannel(router.Channel("errors"))
//	go consumeChannel(router.Channel("default"))
//
//	for entry := range source {
//		router.Dispatch(entry)
//	}
package routing
