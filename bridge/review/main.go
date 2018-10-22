package main

import "fmt"

func bridge(done <-chan interface{}, chs <-chan <-chan interface{}) <-chan interface{} {
	valChan := make(chan interface{})
	go func() {
		defer close(valChan)
		for {
			var stream <-chan interface{}
			select {
			case <-done:
				return
			case v, ok := <-chs:
				if !ok {
					return
				}
				stream = v
			}

			for x := range ordone(done, stream) {
				select {
				case <-done:
					return
				case valChan <- x:
				}
			}
		}
	}()
	return valChan
}

func ordone(done <-chan interface{},inCh <-chan interface{}) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case v, ok := <-inCh:
				if !ok {
					return
				}

				select {
				case outChan <- v:
				case <-done:
				}
			}
		}
	}()
	return outChan
}

func main() {
	genChan := func() <-chan <-chan interface{} {
		outChan := make(chan (<-chan interface{}))
		go func() {
			defer close(outChan)
			for i := 0; i < 50; i++ {
				stream := make(chan interface{}, 2)
				for j := 0; j < 2; j++ {
					stream <- j+i
				}
				close(stream)
				outChan <- stream
			}
		}()
		return outChan
	}
	chs := genChan()
	for v := range bridge(nil, chs) {
		fmt.Printf("%v\n", v)
	}
}
