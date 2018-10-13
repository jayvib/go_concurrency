package main

import (
	"context"
	"math/rand"
	"time"
	"fmt"
	"sync"
)

func intGenerator(ctx context.Context) <-chan int {
	rand.Seed(time.Now().UnixNano())
	intChan := make(chan int)
	go func() {
		defer close(intChan)
		for {
			time.Sleep(500 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			case intChan <- rand.Intn(50):
			}
		}
	}()
	return intChan
}

func tee(ctx context.Context, intChan <-chan int) (_, _, _ <-chan int) {
	out1, out2, out3 := make(chan int), make(chan int), make(chan int)
	go func() {
		defer func() {
			close(out1)
			close(out2)
			close(out3)
		}()
		for v := range intChan {
			var out1, out2, out3 = out1, out2, out3 // create a new copy
			for i := 0; i < 3; i++ {
				select {
				case <-ctx.Done():
					return
				case out1 <- v:
					out1 = nil	// block the channel forever
				case out2 <- v:
					out2 = nil
				case out3 <- v:
					out3 = nil
				}
			}
		}
	}()
	return out1, out2, out3
}

func addBy(x int) func(y int) int {
	return func(y int) int {
		return x + y
	}
}

func multiplyBy(x int) func(y int) int {
	return func(y int) int {
		return x + y
	}
}

func subtractBy(x int) func(y int) int {
	return func(y int) int {
		return x - y
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	defer cancel()

	gen := intGenerator(ctx)
	out1, out2, out3 := tee(ctx, gen)

	wg.Add(3)
	go func(wg *sync.WaitGroup, intChan <-chan int) {
		defer wg.Done()
		fmt.Println("Add process started")
		time.Sleep(500 * time.Millisecond)
		defer fmt.Println("Add process ended")
		addByTwo := addBy(2)
		for v := range intChan {
			fmt.Printf("Added by two: %d\n", addByTwo(v))
		}
	}(&wg, out1)

	go func(wg *sync.WaitGroup, intChan <-chan int) {
		defer wg.Done()
		fmt.Println("Multiply process started")
		time.Sleep(500 * time.Millisecond)
		defer fmt.Println("Multiply process ended")
		multiplyByTwo := multiplyBy(2)
		for v := range intChan {
			fmt.Printf("Multipled by two: %d\n", multiplyByTwo(v))
		}
	}(&wg, out2)

	go func(wg *sync.WaitGroup, intChan <-chan int) {
		defer wg.Done()
		fmt.Println("Subtract process started")
		time.Sleep(500 * time.Millisecond)
		defer fmt.Println("Subtract process ended")
		subtractByTwo := subtractBy(2)
		for v := range intChan {
			fmt.Printf("Subtracted by two: %d\n", subtractByTwo(v))
		}
	}(&wg, out3)

	wg.Wait()
}