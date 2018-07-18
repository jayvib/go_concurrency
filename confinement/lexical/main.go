package main

import "fmt"

func main() {
	chanOwner := func(size int) <-chan int {
		results := make(chan int, size) // it confines the write aspect of this channel to provent goroutines from writing to it.

		go func() {
			defer close(results)
			for i := 0; i < size; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(ch <-chan int) {
		for v := range ch {
			fmt.Println("Receive:", v)
		}
	}

	ch := chanOwner(10)
	consumer(ch)
}
