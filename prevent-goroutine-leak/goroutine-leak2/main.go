package main

import (
	"fmt"
	"math/rand"
	"time"
)

func goroutineLeak() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited")
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()
		return randStream
	}
	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i ++ {
		fmt.Printf("%d: %d\n", i, <-randStream)	// upon exit.. the randStream channel upstream will be block
														// that will cause a goroutine leak.
	}
}

func goroutineLeakSolution() {
	newRandStream := func(done <-chan struct{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done: // we have to tell the goroutine to stop explicitly
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan struct{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i ++ {
		fmt.Printf("%d: %d\n", i, <-randStream)	// upon exit.. the randStream channel upstream will be block
		// that will cause a goroutine leak.
	}
	close(done)
	time.Sleep(1 * time.Second)
	fmt.Println("Done.")
}
