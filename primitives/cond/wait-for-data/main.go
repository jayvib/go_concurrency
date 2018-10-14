package main

import (
	"fmt"
	"sync"
	"time"
)

type Record struct {
	sync.Mutex // guards data
	data string

	cond *sync.Cond
}

func NewRecord() *Record {
	r := Record{}
	r.cond = sync.NewCond(&r)
	return &r
}

func main() {
	var wg sync.WaitGroup

	rec := NewRecord()
	wg.Add(1)

	go func(rec *Record) {
		defer wg.Done() // decrement semaphore counter
		rec.Lock() // enter the critical section
		fmt.Println("gonna wait......ohhh just gonna sleep")
		rec.cond.Wait() // it will unlock..then locks it again only when it wakes up by other go routine.
						// waiting go-routine get suspended(means Go scheduler can execute some other go-routine)

		rec.Unlock() // free up
		fmt.Println("Data: ", rec.data)
		return
	}(rec)

	time.Sleep(2 * time.Second)
	rec.Lock() // enter the critical section
	rec.data = "gopher" // set the value
	rec.Unlock() // frees up
	fmt.Println("waking up the goroutine...")
	rec.cond.Signal() // signal the waiting goroutine.
	  			      // notifies only one of the goroutines(longes waiting goroutine) that are waiting on the
	  			      // condition variable.

	wg.Wait()
}
