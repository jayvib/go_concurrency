package main

import "fmt"

// The unidirectional channel
// <-in send-only; only has the permission to close
// in<- receive-only; don't have a permission to close the channel

func counter(out chan<- int) {
	for x := 1; x <= 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for x := range in {
		out <- x*x
	}
	close(out)
}

func printer(in <-chan int) {
	for x := range in {
		fmt.Println(x)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go counter(naturals)
	go squarer(squares, naturals)

	printer(squares)
}
