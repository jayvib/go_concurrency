package main

import (
	"fmt"
	"time"
	"math/rand"
	"context"
	"sync"
)

// confinement is an exercise for using the confinement concept of goroutine.
// confinement is an idea that whoever make the channel is also the one responsible
// closing it.

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func numberGen(ctx context.Context, n int) <-chan int {
	intChan := make(chan int)

	go func() {
		defer func() {
			fmt.Println("closing the channel...")
			close(intChan)
		}()

		for {
			select {
			case intChan <- rand.Intn(n):
			case <-ctx.Done():
				fmt.Println("Times up!")
				return
			}
			time.Sleep(500 * time.Millisecond) // simulate work
		}
	}()
	return intChan
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println("Value received:", v)
	}
}

// Use Waitgroup when the consumer run in a goroutine.
func consumerWithWaitGroup(wg *sync.WaitGroup, ch <-chan int, name string) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("%s received: %d\n", name, v)
	}
}

func multiplyBy(x int) func(intCh <-chan int, name string) (<-chan string) {
	return func(intCh <-chan int, name string) <-chan string {
		resultChan := make(chan string)
		go func() {
			defer close(resultChan)
			for v := range intCh {
				res := fmt.Sprintf("%s multiplied number %d by %d: %d\n", name, v, x, v*x)
				resultChan <- res
			}
		}()
		return resultChan
	}
}

func main() {
	var wg sync.WaitGroup; _ = wg
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	defer cancel()
	intChan := numberGen(ctx, 100)

	// Example using multiple consumer
	//wg.Add(2)
	//// multiple consumer with one producer
	//go consumerWithWaitGroup(&wg, intChan, "Mars")
	//go consumerWithWaitGroup(&wg, intChan, "Venus")

	// Example using the function currying and closure for having a generic function muliplier. This is an example
	// of fan-in-fan-out-fain-in example.
	multiplyFn := multiplyBy(10)
	resChan1 := multiplyFn(intChan, "Mars")
	resChan2 := multiplyFn(intChan, "Jupiter")

	// using for-select for collecting the results
	resChan1Done := false
	resChan2Done := false
	for !resChan1Done || !resChan2Done {
		select {
		case v, ok := <-resChan1:
			if !ok {
				fmt.Println("resChan1 done...")
				resChan1Done = true
				resChan1 = nil // so that the next loop it will be block forever.
			}
			fmt.Print(v)
		case v, ok := <-resChan2:
			if !ok {
				fmt.Println("resChan2 done...")
				resChan2Done = true
				resChan2 = nil // so that the next loop it will be block forever.
			}
			fmt.Print(v)
		}
	}

	//wg.Wait()
	fmt.Println("Exiting main")
}
