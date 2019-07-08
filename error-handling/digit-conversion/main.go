package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type request struct {
	value string
}

type response struct {
	value int
	err   error
}

func stringToIntConverter(ctx context.Context, reqChan <-chan request) <-chan response {
	respChan := make(chan response)
	go func() {
		defer close(respChan)
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-reqChan:
				time.Sleep(2 * time.Second)
				i, err := strconv.Atoi(req.value)
				respChan <- response{
					value: i,
					err:   err,
				}
			}
		}
	}()
	return respChan
}

func generator(canceler context.CancelFunc) <-chan request {
	reqChan := make(chan request)
	go func() {
		defer func() {
			close(reqChan)
			canceler()
		}()
		reqChan <- request{value: "1"}
		reqChan <- request{value: "2"}
		reqChan <- request{value: "3"}
		reqChan <- request{value: "not valid"}
		reqChan <- request{value: "4"}
		reqChan <- request{value: "not valid"}
	}()
	return reqChan
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reqChan := generator(cancel)
	resChan := stringToIntConverter(ctx, reqChan)
	for response := range resChan {
		if response.err != nil {
			fmt.Println("error:", response.err)
			continue
		}
		fmt.Println(response.value)
	}
}
