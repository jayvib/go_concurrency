package main

import (
	"fmt"
	"sync"
	"time"
)

// This is an attempt to review my understanding
// about conditional variable

type Counter struct {
	sync.Mutex
	count int
	done  bool
	cond  *sync.Cond
}

func NewCounter() *Counter {
	c := &Counter{}
	c.cond = sync.NewCond(c)
	return c
}

func main() {

	updaterFunc := func(wg *sync.WaitGroup, id int, counter *Counter) {
		defer func() {
			wg.Done()
			fmt.Printf("[%d] exiting\n", id)
		}()
		for {
			counter.Lock()
			if counter.done {
				counter.Unlock()
				break
			}
			fmt.Printf("[%d] Waiting for the signal\n", id)
			counter.cond.Wait() // this will unlock... then if theres a signal triggered... it will do Locking automatically
			counter.count++
			fmt.Printf("[%d] Current count value %d\n", id, counter.count)
			counter.Unlock()
		}

	}

	broadcasterFunc := func(wg *sync.WaitGroup, counter *Counter) {
		defer func() {
			wg.Done()
			fmt.Println("Broadcaster func exiting")
		}()
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("Broadcasting...")
			counter.cond.Broadcast()
		}

		counter.done = true
	}

	counter := NewCounter()
	var wg sync.WaitGroup

	updaterNo := 10
	for i := 0; i < updaterNo; i++ {
		wg.Add(1)
		go updaterFunc(&wg, i, counter)
	}

	wg.Add(1)
	go broadcasterFunc(&wg, counter)
	wg.Wait()

}
