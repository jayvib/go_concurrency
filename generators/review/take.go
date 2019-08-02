package main

import "context"

// take take's num values from the value stream and send it
// to the user.
func take(ctx context.Context, valueStream <-chan interface{}, num int) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done:
				return
			case out <- <-valueStream: // get the value from the valueStream and send it to out stream
			}
		}
	}()
	return out
}
