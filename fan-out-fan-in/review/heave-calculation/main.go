package main

import (
	"context"
	"math/rand"
	"time"
)

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

func Repeat(ctx context.Context, values ...interface{}) <-chan interface{} {
	return repeat(ctx, values)
}
func RepeatFn(ctx context.Context, fn func() interface{}) <-chan interface{} {
	return repeatFn(ctx, fn)
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

func Take(ctx context.Context, valueStream <-chan interface{}, num int) <-chan interface{} {
	return take(ctx, valueStream, num)
}

// A list of integers -> Generator
// heavy calculation for every integers -> stage function

func intGenerator() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000000)

}

func main() {
}
