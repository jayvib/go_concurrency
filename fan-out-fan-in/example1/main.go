package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	producer := func(ctx context.Context, id int) <-chan int {
		outChan := make(chan int)
		fmt.Printf("Producer %d is up!\n", id)
		defer fmt.Printf("Producer %d going out\n", id)
		go func() {
			defer close(outChan)
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Producer ID %d is done.\n", id)
					return
				case outChan <- id+rand.Intn(100):
				}
			}
		}()
		return outChan
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	defer cancel()
	prod1 := producer(ctx, 1)
	prod2 := producer(ctx, 2)

	chan1Done := false
	chan2Done := false
	for !chan1Done || !chan2Done {
		select {
		case v, ok := <-prod1:
			if !ok {
				prod1 = nil
				chan1Done = true
				continue
			}
			fmt.Printf("Producer ID %d produce: %d\n", 1, v)
		case x, ok := <-prod2:
			if !ok {
				prod2 = nil
				chan2Done = true
				continue
			}
			fmt.Printf("Producer ID %d produce: %d\n", 2, x)
		}
	}
}
