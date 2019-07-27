package main

type Subscriber interface {
	Notify(interface{}) error
	Topic() string
	Close()
}
