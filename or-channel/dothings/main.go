package main

import (
	"time"
	"fmt"
	"github.com/go_concurrency/or-channel/ordone"
)

func main() {
	dothing := func(name, do string, take time.Duration) <-chan interface{} {
		timer := time.NewTimer(take)
		done := make(chan interface{})

		go func() {
			defer close(done)
			fmt.Printf("%s is doing %s and it take %v to finish.\n",
				name, do, take)
				<-timer.C
		}()

		return done
	}

	start := time.Now()
	<-ordone.OrDone(
		dothing("Foo", "Washing", 5 * time.Minute),
		dothing("Bar", "Cleaning", 8 * time.Hour),
		dothing("Jayson", "Eating chocolate", 5 * time.Second),
	)

	fmt.Printf("Take: %v\n", time.Since(start))
}