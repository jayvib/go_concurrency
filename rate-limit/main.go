package main

import (
	"log"
	"os"
	"github.com/go_concurrency/rate-limit/api"
	"sync"
	"context"
)

func main() {
	defer log.Printf("Done.")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := api.Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("connot Readfile: %v\n", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("connot ResolveAddress: %v\n", err)
			}
			log.Printf("ResolveAddress")
		}()

	}

	wg.Wait()
}