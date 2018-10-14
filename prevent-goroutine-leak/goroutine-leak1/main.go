package main

import (
	"fmt"
	"time"
)

func goroutineLeak() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{}) // when this wasn't close.. it will not garbage collected
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings { // when the strings value is nil... this will be block forever.
									// the goroutine containing dowork will remain in memory for the lifetime of this process.
				completed <- s
			}
		}()
		return completed
	}
	doWork(nil) // a nil goroutine will block forever
	fmt.Println("Done")
}

func goroutineLeakSolution() {
	// as convention done must be in the first paramter of the function.
	doWork := func(done <-chan struct{}, strings <-chan string) <-chan interface{} { // 1
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done: // 2, checking if the done channel has been signaled.
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan struct{})
	terminated := doWork(done, nil)
	go func() { // 3
		// this goroutine will be responsible for signaling the goroutine that the work
		// must be stopped.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done) // sending a signal to done channel
	}()
	<-terminated // 4
	fmt.Println("done")
}

func main() {
	goroutineLeak()
}
