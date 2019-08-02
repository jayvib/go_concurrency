package main

import "context"

// repeat lift the values to the channel until the values
// has been sink then it will reiterate the values and send
// again to the channel. This will keep looping until
// someone call context done.
func repeat(ctx context.Context, values ...interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			for _, v := range values {
				select {
				case <-ctx.Done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func repeatFn(ctx context.Context, fn func() interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done:
				return
			case out <- fn():
			}
		}
	}()
	return out
}
