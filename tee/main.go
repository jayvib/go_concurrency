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
			defer close(out1)
			defer close(out2)
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

	//repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	//	valueStream := make(chan interface{})
	//	go func() {
	//		defer close(valueStream)
	//		for {
	//			for _, v := range values {
	//				select {
	//				case <-done:
	//					return
	//				case valueStream <- v:
	//				}
	//			}
	//		}
	//	}()
	//	return valueStream
	//}
	//
	//// take will send a x-number of data to send into the final stream.
	//take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	//	takeStream := make(chan interface{})
	//	go func() {
	//		defer close(takeStream)
	//		for i := 0; i < num; i++ {
	//			select {
	//			case <-done:
	//				return
	//			case takeStream <- <- valueStream: // <- <- is automatic data casting.
	//			}
	//		}
	//	}()
	//	return takeStream
	//}

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

	out1, out2, out3 := tee(done, gen)

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v out3: %v\n", val1, <-out2, <-out3)
	}
}
