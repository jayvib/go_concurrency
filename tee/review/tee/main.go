package main

import "context"

func ordone(ctx context.Context, data chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-data:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func tee(ctx context.Context, data chan interface{}) (_, _ chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer func() {
			close(out1)
			close(out2)
		}()
		for v := range ordone(ctx, data) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
					return
				case out1 <- v:
					out1 = nil
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func main() {
}
