package main

import (
	"go_concurrency/ordone"
	"fmt"
	"sync"
	"time"
	"math/rand"
)

func bridge(done <-chan interface{}, channels <-chan (<-chan interface{})) <-chan interface{} {
	valstream := make(chan interface{})


	go func() {
		var wg sync.WaitGroup
		rand.Seed(time.Now().UnixNano()) // to provide different random value every run
		defer func() {
			wg.Wait() // wait for its children before exiting.
			fmt.Println("Closing")
			close(valstream)
		}()
		for {
			var stream <-chan interface{}
			select {
			case <-done:
				return
			case v, ok := <-channels:
				if !ok {
					return
				}
				stream = v
			}

			wg.Add(1)
			//fmt.Println("adding worker")
			// so that multiple workers will put value to the stream.
			go func(s <-chan interface{}) {
				defer wg.Done()
				for v := range ordone.OrDone(done, s) {
					time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
					select {
					case valstream <- v:
					case <-done:
					}
				}
			}(stream)
		}
	}()
	return valstream
}

func main() {
	genvals := func() <-chan (<-chan interface{}) {
		chanstream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanstream)
			for i := 0; i < 50; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanstream <- stream
			}
		}()
		return chanstream
	}

	done := make(chan interface{})
	gen := genvals()

	for v := range bridge(done, gen) {
		fmt.Printf("%v ", v)
	}

}
