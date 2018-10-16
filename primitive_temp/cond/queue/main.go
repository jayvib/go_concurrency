package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Let’s expand on this example and show both sides of the equation: a goroutine that is
waiting for a signal, and a goroutine that is sending signals. Say we have a queue of
fixed length 2, and 10 items we want to push onto the queue. We want to enqueue
items as soon as there is room, so we want to be notified as soon as there’s room in
the queue.
 */


func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock() // enter the critical section for the condition so we can modify data pertinent to the condition.
		queue = queue[1:] // 9. simulate dequeuing an item by reassigning the head of the slice to the second item.
		fmt.Println("Removed from queue")
		c.L.Unlock() // 10 exit the critical section since we've successfully dequeued an item.
		c.Signal() // 11 here we let a goroutine waiting on the condition know that something has occured
	}

	for i := 0; i < 10; i++ {
		c.L.Lock() // 3 we enter the critical section for the condition by calling lock on the conditions Locker
		for len(queue) == 2 { // 4 its like telling the the process that... "hey the items in the queue is only 2,...wait for more data to be put to the queue.
			c.Wait() // 5 wait for more data to come in... this will be signaled by different goroutine. It will suspend the main goroutine until a signal on the condition has been sent.
		}
		fmt.Println("adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second) // 6 create a new goroutine that will dequeue an element after one second.
		c.L.Unlock() // 7
	}
}