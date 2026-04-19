package gate

import "context"

// Run reads entries from in, forwards those that pass the gate to out,
// and closes out when in is drained or ctx is cancelled.
func (g *Gate) Run(ctx context.Context, in <-chan map[string]any, out chan<- map[string]any) {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-in:
			if !ok {
				return
			}
			if g.Allow(entry) {
				select {
				case out <- entry:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
