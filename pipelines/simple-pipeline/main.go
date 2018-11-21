package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		defer close(naturals)
		for i := 1;i < 100;i++ {
			naturals <- i
		}
	}()

	go func() {
		defer close(squares)
		for {
			x, ok := <-naturals
			if !ok {
				break
			}
			squares <- x*x
		}
	}()

	for {
		v, ok := <-naturals
		if !ok {
			break
		}
		fmt.Println(v)
	}
}

