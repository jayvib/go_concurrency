package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// This code will attempt to understand what RWMutex is...
//
// Reference: https://hackernoon.com/dancing-with-go-s-mutexes-92407ae927bf
// Use the RWMutext when you can absolutely guarantee that your code
// within your critical section does not mutate shared state.

func main() {
	data := make(map[int]int)
	rand.Seed(time.Now().UnixNano())

	writer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			l.Lock()
			data[rand.Int()] = rand.Int()
			l.ULock()
		}
	}

	readers := func(wg *sync.WaitGroup, l *sync.RWMutex) {
		defer wg.Done()
		l.RLock()
		defer l.RUnlock()
		for k, v := range data {
			fmt.Println(k, v)
		}
	}
	_, _ = writer, reader
}
