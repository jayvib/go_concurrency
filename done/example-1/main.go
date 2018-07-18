package example_1

import (
	"fmt"
	"time"
)

func main() {
	dowork := func(done <-chan struct{}, strings <-chan string, id string) (<-chan struct{}) {
		completed := make(chan struct{})

		go func() {
			defer close(completed)
			defer fmt.Printf("%s exited.\n", id)
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
	completed := dowork(done, nil, "worker 1")
	completed2 := dowork(done, nil, "worker 2")
	completed3 := dowork(done, nil, "worker 3")

	go func() {
		timer := time.NewTimer(5 * time.Second)
		fmt.Println("Timer start.")
		<-timer.C
		fmt.Println("Times up! You're too slow! Cancelling work...")
		close(done) // it's another phrase to say na "ayaw nala ipadayun it trabaho mahinay kaman" hehehe.
	}()

	<-completed
	<-completed2
	<-completed3
}
