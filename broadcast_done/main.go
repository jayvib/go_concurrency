package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)

	// launch goroutine they will wait/block
	for i := 0; i < 5; i++ {
		go func(i int) {
			_, ok := <-c
			fmt.Printf("closed %d, %t\n", i, ok)
		}(i)
	}

	close(c)
	time.Sleep(time.Second)
}
