package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	doWork := func(done <-chan interface{}, intSlice ...int) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1)
		resultStream := make(chan int)

		go func() {
			defer func() {
				close(heartbeatStream)
				close(resultStream)
			}()

			duration := time.Duration(rand.Intn(10))

			fmt.Printf("Initializing objects... ready after %d seconds\n", duration)
			time.Sleep(duration * time.Second)

			fmt.Println("Database connected")
			time.Sleep(1 * time.Second)
			fmt.Println("Elasticsearch connected")
			time.Sleep(1 * time.Second)
			fmt.Println("Redis conneced")
			time.Sleep(1 * time.Second)

			workgen := time.Tick(200 * time.Millisecond)

			for _, n := range intSlice {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}

				select {
				case <-done:
					return
				case <-workgen:
					resultStream <- n
				}
			}

		}()
		return heartbeatStream, resultStream
	}

	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 4}

	heartbeat, results := doWork(done, intSlice...)

	<-heartbeat

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			fmt.Printf("index %v: expected %v, but received %v", i, expected, r)
		}
		i++
		fmt.Println(r)
	}
	fmt.Println("Success")
}
