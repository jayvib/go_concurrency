package main

import (
	"context"
	"time"
	"math/rand"
	"fmt"
)

type result struct {
	num int
	res bool
	err error
}

type filterFunc func(int) bool

func (fn filterFunc) run(ctx context.Context, intChan <-chan int) <-chan result {
	out := make(chan result)

	go func() {
		defer close(out)

		for val := range intChan {
			var r result
			r.num = val
			select {
			case <-ctx.Done():
				r.err = ctx.Err()
				out <- r
				return
			default:
				r.res = fn(val)
				out <- r
			}
		}
	}()
	return out
}

func NewFilterFunc(fn func(int) bool ) filterFunc {
	return filterFunc(fn)
}

func NewIsGreaterFunc(compareFrom int) filterFunc {
	return filterFunc(func(n int) bool {
		return compareFrom > n
	})
}

func numgenerator(ctx context.Context, maxnum int) <-chan int {
	intChan := make(chan int)
	go func() {
		defer close(intChan)

		for {
			select {
			case <-ctx.Done():
				return
			case intChan <- rand.Intn(maxnum):
			}
		}
	}()
	return intChan
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	isGreater := NewIsGreaterFunc(10)
	run := isGreater.run

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	intChan := numgenerator(ctx, 50)

	for res := range run(ctx, intChan) {
		if res.err != nil {
			fmt.Println("got an error:", res.err.Error())
			continue
		}
		fmt.Printf("Number: %d CompareFrom: %d Result: %t\n", res.num, 10, res.res)
	}
}
