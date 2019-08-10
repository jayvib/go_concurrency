package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

// bridge is completely independent without relying the lifetime of the function
// from the channels owned by other goroutine.
//
// Basically, bridge function will just accepting the read-only channels and
// read data from it and send it to the downstream channel.
//
// when there's a problem in the upstream operation the bridge will not be affected
// removing the possibility of recreating or recalling another bridge function call.
func bridge(ctx context.Context, chs <-chan <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ch, ok := <-chs:
				if !ok {
					return
				}

				for v := range ordone(ctx, ch) {
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}
	}()
	return out
}

func ordone(ctx context.Context, ch <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
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

func generator() <-chan (<-chan interface{}) {
	outchan := make(chan (<-chan interface{}))
	go func() {
		defer close(outchan)
		for i := 0; i < 100; i++ {
			stream := make(chan interface{}, 5)
			for j := 0; j < 5; j++ {
				stream <- j + i
			}
			close(stream)
			outchan <- stream
		}
	}()
	return outchan
}

func main() {
	genChan := generator()
	for v := range bridge(context.Background(), genChan) {
		logrus.Printf("%v\n", v)
	}
}
