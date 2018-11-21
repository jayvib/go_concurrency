package main

import "sync"

func main() {
	fanIn := func(done <-chan interface{}, chans ...<-chan interface{}) <-chan interface{} { //1
		outChan := make(chan interface{}) // to combine all the outputs
		var wg sync.WaitGroup // four goroutines that will drain the input from the upstream operation. 2
							  // co-ordinate goroutines.
		multiplexChan := func(wg *sync.WaitGroup, inChan <-chan interface{}) { //3
			defer wg.Done() // signal the waitgroup that this goroutine is done
			for v := range inChan {
				select {
				case <-done:
					return
				case outChan <- v:
				}
			}
		}

		wg.Add(len(chans)) // 4
		for _, ch := range chans {
			go multiplexChan(&wg, ch)
		}

		go func() { //5
			wg.Wait() // wait all goroutines to finish before closing the out channel
			close(outChan)
		}()
		return outChan
	}

}
