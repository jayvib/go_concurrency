package main

import (
	"fmt"
	"time"
)

// this function enables you to combine any number of channels together into
// a single channel that will close as soon as any of tis component channels are closed, or written to.

// a signle done channel that closes if any of its component channels close.
func Or(dones ...<-chan interface{}) <-chan interface{} { // 1
	switch len(dones) {
	case 0: // 2
		return nil
	case 1: // 3
		return dones[0]
	}

	orDone := make(chan interface{})
	go func() { // 4
		// this goroutine will send a signal to orDone channel if
		// any of the done dones has been signaled to stop.
		defer close(orDone)
		switch len(dones) {
		case 2: //5
			select {
			case <-dones[0]:
			case <-dones[1]:
			}
		default: // 6
			select {
			case <-dones[0]:
			case <-dones[1]:
			case <-dones[2]:
			case <-Or(append(dones[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		done := make(chan interface{})
		go func() {
			defer close(done)
			time.Sleep(after)
		}()
		return done
	}

	start := time.Now()
	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
