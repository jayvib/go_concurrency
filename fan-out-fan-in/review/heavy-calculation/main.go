package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// repeat lift the values to the channel until the values
// has been sink then it will reiterate the values and send
// again to the channel. This will keep looping until
// someone call context done.
func repeat(ctx context.Context, values ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for {
			for _, v := range values {
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func repeatFn(ctx context.Context, fn func() int) <-chan int {
	out := make(chan int)

	go func() {
		defer func() {
			close(out)
			logrus.Println("RepeatFn exiting")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case out <- fn():
			}
		}
	}()
	return out
}

// take take's num values from the value stream and send it
// to the user.
func take(ctx context.Context, valueStream <-chan int, num int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- <-valueStream: // get the value from the valueStream and send it to out stream
			}
		}
	}()
	return out
}

func Take(ctx context.Context, valueStream <-chan int, num int) <-chan int {
	return take(ctx, valueStream, num)
}

func Repeat(ctx context.Context, values ...int) <-chan int {
	return repeat(ctx, values...)
}
func RepeatFn(ctx context.Context, fn func() int) <-chan int {
	return repeatFn(ctx, fn)
}

// A list of integers -> Generator
// heavy calculation for every integers -> stage function

func intGenerator() int {
	return rand.Intn(1000000)

}

func calculate(ctx context.Context, ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range ordone(ctx, ch) {
			// simulate heavy calculation
			duration := time.Duration(rand.Intn(5))
			logrus.Println("Calculating...")
			time.Sleep(duration * time.Second)
			out <- v * rand.Intn(50)
		}
	}()
	return out
}

// fanin multiplex the output from the workers and combine
// it into a single channel
func fanIn(ctx context.Context, chs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedChan := make(chan int)

	multiplexFunc := func(ctx context.Context, wg *sync.WaitGroup, ch <-chan int) {
		defer wg.Done()
		for v := range ordone(ctx, ch) {
			multiplexedChan <- v
		}
	}

	for _, ch := range chs {
		wg.Add(1)
		go multiplexFunc(ctx, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(multiplexedChan)
	}()

	return multiplexedChan
}

func ordone(ctx context.Context, ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch: // check if the ok, meaning the channel is not yet closed
				if !ok {
					return
				}

				// sometimes when the calculation is heavy
				// this goroutine can't know that
				// the upstream can cancel the work at
				// any time.
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	logrus.SetFormatter(new(logrus.TextFormatter))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	intChan := RepeatFn(ctx, intGenerator)
	takeIntChan := Take(ctx, intChan, 20)

	for v := range ordone(ctx, takeIntChan) {
		logrus.Println(v)
	}

	//intChans := make([]<-chan int, 0)
	//workers := 3
	//for i := 0; i < workers; i++ {
	//	intChans = append(intChans, calculate(ctx, takeIntChan))
	//}

	//var resultChan <-chan int

	//go func() {
	//	resultChan = fanIn(ctx, intChans...)
	//}()

	//for v := range resultChan {
	//	logrus.Println("Result:", v)
	//}
}
