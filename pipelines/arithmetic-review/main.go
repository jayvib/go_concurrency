package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// This program is an attempt to review the
// pipeline pattern...
//

func main() {
	generator := func(ctx context.Context) <-chan int {
		rand.Seed(time.Now().UnixNano())
		intChan := make(chan int)
		go func() {
			defer close(intChan)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(1 * time.Second)
					intChan <- rand.Intn(10)
				}
			}
		}()
		return intChan
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	intChan := generator(ctx)
	printInt(multiply(ctx, add(ctx, intChan, 100), 2))
}

func multiply(ctx context.Context, intChan <-chan int, multiplier int) <-chan int {
	intChanOut := make(chan int)
	go func() {
		defer close(intChanOut)
		for i := range intChan {
			select {
			case <-ctx.Done():
				return
			default:
				intChanOut <- (i * multiplier)
			}
		}
	}()
	return intChanOut
}

func add(ctx context.Context, intChan <-chan int, adder int) <-chan int {
	intChanOut := make(chan int)
	go func() {
		defer close(intChanOut)
		for i := range intChan {
			select {
			case <-ctx.Done():
				return
			default:
				intChanOut <- i + adder
			}
		}
	}()
	return intChanOut
}

func printInt(ch <-chan int) {
	for i := range ch {
		fmt.Println(i)
	}
}
