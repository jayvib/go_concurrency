package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func Generator(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case outChan <- fn():
			}
		}
	}()
	return outChan
}

func Take(done <-chan interface{}, inChan <-chan interface{}, num int) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case outChan <- <- inChan:
			}
		}
	}()
	return outChan
}

func ToInt(done <-chan interface{}, inChan <-chan interface{}) <-chan int {
	outChan := make(chan int)
	go func() {
		defer close(outChan)
		for i := range inChan {
			select {
			case <-done:
				return
			case outChan <- i.(int):
			}
		}
	}()
	return outChan
}

func primeFinder(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})
	go func() {
		defer close(primeStream)
		for integer := range intStream {
			integer -= 1
			prime := true
			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()
	return primeStream
}

func main() {
	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		// multiplex will wait for the input from c then send it to multiplxedStream channel.
		multiplex := func(c <-chan interface{}) {
			defer wg.Done() // make sure that the semaphore counter will be decremented
			for v := range c { // loop over the values from channel
				select {
				case <-done: // when closed this function will exit
					return
				case multiplexedStream <- v:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()
		return multiplexedStream
	}

	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	random := func() interface{} {
		return rand.Intn(5000000)
	}
	randIntStream := ToInt(done, Generator(done, random)) // pipeline pattern
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Printf("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range Take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))
}
