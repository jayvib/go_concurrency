package main

import (
	"go_concurrency/ordone"
	"fmt"
	"time"
	"math/rand"
)

func main() {
	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		out3 := make(chan interface{})
		go func() {
			defer func() {
				close(out1)
				close(out2)
				close(out3)
			}()
			for val := range ordone.OrDone(done, in) {
				var out1, out2, out3 = out1, out2, out3 // new copy
				for i := 0; i < 3; i++ {
					select {
					case <-done:
						return
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					case out3 <- val:
						out3 = nil
					}
				}
			}
		}()
		return out1, out2, out3
	}

	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	// take will send a x-number of data to send into the final stream.
	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <- valueStream: // <- <- is automatic data casting.
				}
			}
		}()
		return takeStream
	}

	generator := func(done <-chan interface{}) <-chan interface {} {
		out := make(chan interface{})
		rand.Seed(time.Now().UnixNano())
		go func() {
			for {
				select {
				case <-done:
					return
				case out <- rand.Int():
				}
				time.Sleep(500 * time.Millisecond)
			}
		}()
		return out
	}

	done := make(chan interface{})
	defer close(done)

	gen := generator(done)
	_ = gen
	out1, out2, out3 := tee(done, take(done, repeat(done, 1, 2, 3), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v out3: %v\n", val1, <-out2, <-out3)
	}
}

// tee2 takes a done channel and the input channel. It will return 2 goroutines where it has an own
// copy of the value received from inChan.
func tee2(done, inChan <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func(){
		defer func() {
			close(out1)
			close(out2)
		}()
		for v := range ordone.OrDone(done, inChan) {
			var out1, out2 = out1, out2 // every loop the original pointer instance will be copied
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- v:
					out1 = nil // will block in the next iteration.
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}
