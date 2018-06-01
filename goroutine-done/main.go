package main

import (
	"fmt"
	"time"
)

func main() {
	dowork := func(done chan struct{}, strings <-chan string) (<-chan struct{}) {
		completed := make(chan struct{})

		go func() {
			defer fmt.Println("dowork exited.")
			defer close(completed)
			for {
				select {
				case s := <-strings:
					fmt.Println("Receive:", s)
				case <-done:
					return
				}
			}
		}()
		return completed
	}

	done := make(chan struct{})
	completed := dowork(done, nil)

	go func() {
		timer := time.NewTimer(5 * time.Second)
		fmt.Println("Timer start.")
		<-timer.C
		fmt.Println("Times up!")
		done <- struct{}{}
	}()

	<-completed // this will cause deadlock when In not gonna pass a done channel since the strings chan will gonna wait forever.
	fmt.Println("Main program exited.")
}
