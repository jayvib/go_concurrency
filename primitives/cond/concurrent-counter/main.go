package main

import (
	"fmt"
	"sync"
	"time"
)

// concurrent-counter will simulate a multiple updater to a single variable(memory)
// by using the sync.Cond.Broadcast()

type counter struct {
	sync.Mutex
	count int

	cond *sync.Cond
}

func main() {
	exit := false
	countUpdater := func(i int, count *counter) {
		for {
			count.Lock()
			if !exit {
				fmt.Printf("[%d] Waiting for the signal\n", i)
				count.cond.Wait()
				count.count++
				fmt.Printf("[%d] Current value: %d\n", i, count.count)
			}
			count.Unlock()
		}
	}

	notifier := func(c *counter) {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			fmt.Println("broadcasting...")
			c.cond.Broadcast()
		}
		exit = true
	}

	c := &counter{}
	c.cond = sync.NewCond(c)
	var wg sync.WaitGroup
	numUpdater := 10
	wg.Add(numUpdater)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			countUpdater(i, c)
		}(i)
	}

	go notifier(c)
	wg.Wait()
}
