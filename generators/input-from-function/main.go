package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator(done chan struct{}, fn func() interface{}) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer func() {
			fmt.Println("Closing the generator")
			close(outChan)
		}()
		for {
			time.Sleep(500 * time.Millisecond)
			select {
			case <-done:
				return
			case outChan <- fn():
			}
		}
	}()
	return outChan
}

func main() {
	done := make(chan struct{})
	rand.Seed(time.Now().UnixNano())
	dataFunc := func() interface{} {
		return rand.Int()
	}
	genChan := generator(done, dataFunc)

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	for d := range genChan {
		fmt.Println(d)
	}
}
