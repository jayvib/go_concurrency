package main

import (
	"math/rand"
	"time"
	"fmt"
	"go_concurrency/ordone"
)

func main() {
	gen := func(done <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})

		go func() {
			defer close(out)
			rand.Seed(time.Now().UnixNano())

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

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Times up!")
		close(done)
	}()

	ch := gen(done)

	for v := range ordone.OrDone(done, ch) {
		fmt.Println(v)
	}


}
