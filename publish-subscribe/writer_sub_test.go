package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestWriter(t *testing.T) {
	msg := "hello"
	sub := NewWriterSubscriber(0, nil)
	defer sub.Close()
	stdoutPrinter := sub.(*writerSubscriber)

	var wg sync.WaitGroup
	wg.Add(1)
	stdoutPrinter.Writer = &mockWriter{
		testFunc: func(res string) {
			defer wg.Done()
			if !strings.Contains(res, msg) {
				t.Fatalf("Incorrect string: %s", res)
			}
			fmt.Println("ohh yeah")
		},
	}

	err := sub.Notify(msg)
	if err != nil {
		wg.Done()
		t.Error(err)
	}

	wg.Wait()
	sub.Close()
}
