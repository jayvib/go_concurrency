package main

import (
	"fmt"
	"runtime"
)

func Generator(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case outChan <- fn():
			}
		}
	}()
	return outChan
}

func Take(done <-chan interface{}, inChan <-chan interface{}, num int) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case outChan <- inChan:
			}
		}
	}()
	return outChan
}

func ToInt(done <-chan interface{}, inChan <-chan interface{}) <-chan int {
	outChan := make(chan int)
	go func() {
		defer close(outChan)
		for i := range inChan {
			select {
			case <-done:
				return
			case outChan <- i.(int):
			}
		}
	}()
	return outChan
}



func main() {
	fmt.Println(runtime.NumCPU())
}
