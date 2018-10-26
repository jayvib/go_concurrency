package main

import (
	"sync"
	"context"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) // 1, create a context that cancellable
	defer cancel() // do cleanup before exiting so that there will be no goroutine resource leak

	wg.Add(1)
	go func() {
		defer wg.Done() // decrement semaphore
		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel() // 2
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx) // generate a greeting message
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) // 3, create a context that has timeout of 1 second
	defer cancel() // do some cleanup to that uses this context. or will be saying the that this context has cancelled already.
				   // no need for further work.

	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", errors.New("unsupported locale")
}

func locale(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err() // 4
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
