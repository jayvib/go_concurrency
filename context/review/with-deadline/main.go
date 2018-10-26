package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// with-deadline: illustrate how to use the context.WithDeadLine

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10 * time.Second))
	defer cancel()
	msg, err := do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg)
}

func do(ctx context.Context) (string, error) {
	time.Sleep(6 * time.Second)
	if deadline, ok := ctx.Deadline(); ok {
		// check if the deadline is near to a specific point of time
		if deadline.Sub(time.Now().Add(5 * time.Second)) <= 0 { // check if the deadline has 5 or less second left.
			return "", context.DeadlineExceeded
		}
	}
	return "Hello World", nil
}
