package main

import (
	"math/rand"
	"time"
	"fmt"
	"go_concurrency/ordone"
)

func tee(done <-chan interface{}, valStream <-chan interface{}) (_, _, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	out3 := make(chan interface{})

	go func() {
		defer func() {
			close(out1)
			close(out2)
			close(out3)
		}()

		for val := range ordone.OrDone(done, valStream) {
			var out1, out2, out3 = out1, out2, out3
			for i := 0; i < 3; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil // block forever
				case out2 <- val:
					out2 = nil
				case out3 <- val:
					out3 = nil
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return out1, out2, out3
}

func intGenerator(done chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			select {
			case <-done:
				return
			case out <- rand.Int():
			}
		}
	}()
	return out
}

func main() {

	done := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Times-up!")
		close(done)
	}()

	gen := intGenerator(done)
	out1, out2, out3 := tee(done, gen)

	for val1 := range out1 {
		fmt.Printf("out1: %v out2: %v out3: %v\n", val1, <-out2, <-out3)
	}
}
