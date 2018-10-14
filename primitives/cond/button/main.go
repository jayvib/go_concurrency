package main

import (
	"fmt"
	"sync"
	"time"
)

type Button struct { // 1. define a type button that contains a condition.
	Clicked *sync.Cond
}

func main() {
	button := Button{ Clicked: sync.NewCond(&sync.Mutex{}) }

	subscribe := func(c *sync.Cond, fn func()) { // 2
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1) // add a number of goroutine running which is here has 1.
		go func() {
			goroutineRunning.Done()
			c.L.Lock() // lock the critical section
			defer c.L.Unlock() // defer unlock so that whatever happen this critical section will be unlocked
			fmt.Println("waiting")
			c.Wait() // waiting for the signal from the conditional variable
			fn() // execute the function
		}()
		goroutineRunning.Wait() // wait for the goroutine to finish
	}

	var clickRegistered sync.WaitGroup // 3
	clickRegistered.Add(3)
	subscribe(button.Clicked, func(){ // 4
		defer clickRegistered.Done()
		fmt.Println("Displaying annoying dialog box!")
	})

	subscribe(button.Clicked, func(){ // 5
		defer clickRegistered.Done()
		fmt.Println("Mouse clicked.")
	})

	subscribe(button.Clicked, func(){ // 6
		defer clickRegistered.Done()
		fmt.Println("Maximizing window.")
	})
	time.Sleep(2 * time.Second)
	fmt.Println("broadcasting")
	button.Clicked.Broadcast() // 7
	clickRegistered.Wait()
}
