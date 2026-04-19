// Package gate implements a conditional pass-through stage for the logpipe
// processing pipeline.
//
// A Gate is configured with a field name, a comparison operator (eq, neq, gt,
// lt), and a reference value. Only log entries whose named field satisfies the
// comparison are forwarded downstream; all others are silently dropped.
//
// Example usage:
//
//	g, err := gate.New(gate.Config{
//		Field:           "level",
//		Op:              gate.OpEq,
//		Value:           "error",
//		CaseInsensitive: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	g.Run(ctx, in, out)
package gate
