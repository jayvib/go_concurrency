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

	balanceChan := make(chan int)
	tranChan := make(chan bool) // needs to close explicitly

	balance = 1000 // this will cause race condition
	transactionNo = 0
	fmt.Println("Starting balance: $", balance)

	wg.Add(1)
	for i := 0; i < 100; i++ {
		go func(ii int, trChan chan bool, balChan chan int) {
			transactionAmount := rand.Intn(25)
			balChan <- transactionAmount
			transaction(transactionAmount)
			if ii == 99 {
				fmt.Println("Should be quittin time.")
				close(balanceChan)
				trChan <- true
			}
		}(i, tranChan, balanceChan)
	}

	go transaction(0)

	go func() {
		defer close(tranChan)
		defer wg.Done()
		breakpoint := false
		for {
			if breakpoint {
				break
			}

			select {
			case <- tranChan:
				fmt.Println("Transactions finished")
				wg.Done()
			case amount := <-balanceChan:
				fmt.Println("Transaction for $", amount)
				if (balance - amount) < 0 {
					fmt.Println("Transaction failed")
				} else {
					balance -= amount
					fmt.Println("Transaction succeded")
				}
				fmt.Println("Balance now $", balance)
			case isDone := <-tranChan:
				if isDone {
					fmt.Println("done")
					breakpoint = true
				}
			}
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