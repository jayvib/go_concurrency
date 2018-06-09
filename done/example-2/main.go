package main

import (
	"math/rand"
	"fmt"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())

	newRandStream := func(done <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(intStream)
			for {
				select {
				case intStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return intStream
	}

	done := make(chan interface{})
	intStream := newRandStream(done)

	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		fmt.Println("Times up!")
		close(done)
	}()

	for v := range intStream {
		time.Sleep(500 * time.Millisecond) // simulate work
		fmt.Println(v)
	}
}
