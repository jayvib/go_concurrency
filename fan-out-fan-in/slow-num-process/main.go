package main

import (
	"fmt"
	"github.com/jayvib/concutil"
	"github.com/kubernetes/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	done := make(chan struct{})
	randomInt := func() interface {}{
		return rand.Intn(50000)
	}
	intGen := concutil.ToInt(done, concutil.Generator(done, randomInt))
	go func() {
		time.Sleep(20 * time.Second)
		close(done)
	}()

	numProcessors := runtime.NumCPU()
	channels := make([]<-chan int, numProcessors)
	for i := 0; i < numProcessors; i++ {
		channels[i] = processNum(done, intGen)
	}

	for v := range concutil.Take(done, fanIn(done, channels...), 5) {
		fmt.Println(v)
	}
	close(done)
	//go printValue(intGen)
	//time.Sleep(5 * time.Second)
	//close(done)
	//fmt.Println("done")
}


func fanIn(done <-chan struct{}, channels ...<-chan int) <-chan interface{} {
	var wg sync.WaitGroup
	outChan := make(chan interface{})
	multiplexFunc := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			select {
			case <-done:
				return
			case outChan <- v:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplexFunc(c)
	}
	go func() {
		wg.Wait()
		close(outChan)
	}()
	return outChan
}

func processNum(done <-chan struct{}, intCh <-chan int) <-chan int {
	outChan := make(chan int)
	go func() {
		defer close(outChan)
		for v := range intCh {
			time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
			v *= 100
			select {
			case <-done:
				return
			case outChan <- v:
			}
		}
	}()
	return outChan
}

func printValue(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
