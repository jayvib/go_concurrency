package main

import (
	"context"
	"fmt"
)

type Publisher interface {
	Start()
	Subscribe(sub Subscriber)
	Unsubscribe(sub Subscriber)
	PublishingCh(msg interface{})
	Stop()
}

func NewPublisher() Publisher {
	ctx, cancel := context.WithCancel(context.Background())

	return &publisher{
		subscribers: make(map[string]Subscriber),
		in:          make(chan interface{}),
		addSubCh:    make(chan Subscriber),
		removeSubCh: make(chan Subscriber),
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}

// WISDOM: When using channels, we will need a channel for each
// action that can be considered DANGEROUS
type publisher struct {
	subscribers map[string]Subscriber
	addSubCh    chan Subscriber
	removeSubCh chan Subscriber
	in          chan interface{}
	ctx         context.Context // the stop channel must be called when we want to kill all Goroutines
	cancelFunc  func()
}

func (p *publisher) Start() {
	go func() {
		fmt.Println("Publisher started...")
		for {
			// listen from the incoming messages
			select {
			case msg := <-p.in:
				fmt.Println("Message received:", msg)
				for _, sub := range p.subscribers {
					sub.Notify(msg)
				}
			case sub := <-p.addSubCh:
				fmt.Printf("Adding New Subscriber: %v\n", sub.Topic())
				p.subscribers[sub.Topic()] = sub
			case sub := <-p.removeSubCh:
				fmt.Printf("Removing Subscriber: %v\n", sub.Topic())
				candidate := p.subscribers[sub.Topic()]
				fmt.Println("Removing:", candidate.Topic())
				delete(p.subscribers, sub.Topic())
				candidate.Close()
			case <-p.ctx.Done():
				fmt.Println("Stoping the publisher")
				for _, sub := range p.subscribers {
					sub.Close()
					close(p.addSubCh)
					close(p.removeSubCh)
					close(p.in)
				}
				return
			}
		}
	}()
}

func (p *publisher) Subscribe(sub Subscriber) {
	p.addSubCh <- sub
}

func (p *publisher) Unsubscribe(sub Subscriber) {
	p.removeSubCh <- sub
}

func (p *publisher) PublishingCh(msg interface{}) {
	p.in <- msg
}

func (p *publisher) Stop() {
	p.cancelFunc()
}
