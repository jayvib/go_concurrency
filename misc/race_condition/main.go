package main

import (
	"math/rand"
	"time"
	"runtime"
	"sync"
	"fmt"
)

var balance int
var transactionNo int

func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup

	tranChan := make(chan bool) // needs to close explicitly

	balance = 1000 // this will cause race condition
	transactionNo = 0
	fmt.Println("Starting balance: $", balance)

	wg.Add(1)
	for i := 0; i < 100; i++ {
		go func(ii int, trChan chan(bool)) {
			transactionAmount := rand.Intn(25)
			transaction(transactionAmount)
			if ii == 99 {
				trChan <- true
			}
		}(i, tranChan)
	}

	go transaction(0)

	go func() {
		defer close(tranChan)
		select {
		case <- tranChan:
			fmt.Println("Transactions finished")
			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Println("Final balance: $", balance)

}

func transaction(amount int) bool {
	approved := false
	if (balance-amount) < 0 {
		approved = false
	} else {
		approved = true
		balance -= amount
	}
	approvedText := "declined"
	if approved {
		approvedText = "approved"
	}
	transactionNo++
	fmt.Println(transactionNo, "Transaction for $", amount, approvedText)
	fmt.Println("\tRemaining balance $", balance)
	return approved
}