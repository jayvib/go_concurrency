// Pipelines is helpful to break long process into chunk
// of mini-processes this to decompose the implementations
// into function.
package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
)

func subtractBy(x int) func(done <-chan interface{}, intStream <-chan int) <-chan int {
	return func(done <-chan interface{}, intStream <-chan int) <-chan int {
		outChan := make(chan int)
		go func() {
			for v := range intStream {
				select {
				case <-done:
					return
				case outChan <- v-x:
				}
			}
		}()
		return outChan
	}
}

func main() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int, len(integers))
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i*multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)

		go func() {
			defer close(addedStream)
			for n := range intStream {
				select {
				case <-done:
					return
				case addedStream <- n+additive:
				}
			}
		}()
		return addedStream
	}


	done := make(chan interface{})
	c := make(chan os.Signal, 1)

	go func() {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		fmt.Println("exiting")
		defer close(done)
	}()

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := subtractBy(10)(done, multiply(done, add(done, multiply(done, intStream, 2), 1), 2)) // y-combinator

	for v := range pipeline {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(v)
	}

}