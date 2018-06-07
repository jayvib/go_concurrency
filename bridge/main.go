package main

import (
	"go_concurrency/ordone"
	"fmt"
)

func main() {
	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{}) // where the value will be flow out
		go func() {
			defer close(valStream) // deferred close when the function exists
			for { // 2
				var stream <-chan interface{} // a variable where to store the received read-only channel
				select {
				case maybeStream, ok := <-chanStream: // receive the read-only channel then check also if this channel is close
					if !ok {
						return // if close then exit
					}
					stream = maybeStream // save the received read-only channel to the stream.
				case <-done:
					return // if done then exit
				}

				for val := range ordone.OrDone(done, stream) { // use OrDone so that whenere
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream) // a channel that has 1 buffered value that is ready to write but can't be written
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}

}