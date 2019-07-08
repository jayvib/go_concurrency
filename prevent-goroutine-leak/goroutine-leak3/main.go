package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// This program will be a review for handling the goroutine leaks
func generator(ctx context.Context, dcount int) <-chan int {
	rand.Seed(time.Now().UnixNano())
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < dcount; i++ {
			time.Sleep(1 * time.Second) // this will simulate any work
			select {
			case <-ctx.Done(): // this will be called when the time deadline meet
				return
			default:
				ch <- rand.Int() // send int to the channel
			}
		}
	}()
	return ch
}

func consumer(ch <-chan int) {
	fmt.Println("consumer started")
	defer fmt.Println("consumer exiting")
	for i := range ch { // this will iterate data that have send into the channel
		fmt.Println(i)
	}
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(),
		time.Now().Add(5*time.Second))
	defer cancel()
	intChan := generator(ctx, 10)
	consumer(intChan)
}
