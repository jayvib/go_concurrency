package main

import (
	"time"
	"fmt"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} { // 1
		switch len(channels) {
		case 0: // 2
			return nil
		case 1: // 3
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() { // 4
			defer close(orDone)
			switch len(channels) {
			case 2: // 5
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // 6
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): // 6
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2 * time.Hour),
		sig(5 * time.Minute),
		sig(7 * time.Second),
		sig(10 * time.Hour),
		sig(4 * time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
