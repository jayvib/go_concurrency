package main

import (
	"sync"
	"testing"
)

// A mock for the publisher
type mockSubscriber struct {
	notifyTestingFunc func(msg interface{})
	closeTestingFunc  func()
	topicFunc         func() string
}

func (s mockSubscriber) Close() {
	s.closeTestingFunc()
}

func (s mockSubscriber) Notify(msg interface{}) error {
	s.notifyTestingFunc(msg)
	return nil
}

func (s mockSubscriber) Topic() string {
	return s.topicFunc()
}

func TestPublisher(t *testing.T) {
	msg := "Hello"
	pub := NewPublisher()

	var wg sync.WaitGroup
	sub := mockSubscriber{
		notifyTestingFunc: func(m interface{}) {
			v, ok := m.(string)
			if !ok {
				t.Fail()
			}

			if msg != v {
				t.Errorf("Expected result not match: got %s want %s", v, msg)
			}
		},
		closeTestingFunc: func() {
			wg.Done()
		},
		topicFunc: func() string {
			return "testing"
		},
	}

	pub.Start()
	pub.Subscribe(sub)
	wg.Add(1)

	pub.PublishingCh(msg)

	pub.Unsubscribe(sub)
	wg.Wait()
	pub.Stop()
}
