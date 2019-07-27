package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func NewWriterSubscriber(id int, out io.Writer) Subscriber {

	if out == nil {
		out = os.Stdout
	}
	sub := &writerSubscriber{
		in:     make(chan interface{}),
		id:     id,
		Writer: out,
	}

	go func() {
		for msg := range sub.in {
			fmt.Fprintf(sub.Writer, "(W%d): %v\n", sub.id, msg)
		}
	}()
	return sub
}

type writerSubscriber struct {
	in     chan interface{} // the problem is....how the sender knew that this chanel is already close?
	id     int
	Writer io.Writer
	topic  string
}

func (s *writerSubscriber) Notify(msg interface{}) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%#v", rec)
		}
	}()
	select {
	case s.in <- msg: // there might be an error occur for this.... sending a closed channel
	case <-time.After(time.Second):
		err = fmt.Errorf("Timeout")
	}
	return
}

func (s *writerSubscriber) Close() {
	close(s.in)
}

func (s *writerSubscriber) Topic() string {
	return s.topic
}

type mockWriter struct {
	testFunc func(string)
}

func (m *mockWriter) Write(p []byte) (int, error) {
	m.testFunc(string(p))
	return len(p), nil
}
