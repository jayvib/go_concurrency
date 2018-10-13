package main

import "sync"

func main() {
	fanIn := func(done <-chan interface{}, chans ...<-chan interface{}) <-chan interface{} {
		outChan := make(chan interface{})
		var wg sync.WaitGroup

		multiplexChan := func(wg *sync.WaitGroup, inChan <-chan interface{}) {
			defer wg.Done()
			for v := range inChan {
				select {
				case <-done:
					return
				case outChan <- v:
				}
			}
		}

		wg.Add(len(chans))
		for _, ch := range chans {
			go multiplexChan(&wg, ch)
		}

		go func() {
			wg.Wait()
			close(outChan)
		}()
		return outChan
	}


}
