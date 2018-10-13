package main

import (
	"fmt"
	"sync"
)

func main() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeter = 5
	var wg sync.WaitGroup
	wg.Add(numGreeter)

	for i := 1; i <= numGreeter; i++ {
		go hello(&wg, i)
	}
	wg.Wait()
}
