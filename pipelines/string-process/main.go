package main

import (
	"fmt"
	"strings"
)

func ToLower(ch <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for d := range ch {
			out <- strings.ToLower(d)
		}
	}()
	return out
}

func appendLetter2(ch <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for d := range ch {
			out <- d + " 2"
		}
	}()
	return out
}

func AppendLetter(char string) func(ch <-chan string) <-chan string{
	return func(ch <-chan string) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for d := range ch {
				out <- d + " " + char
			}
		}()
		return out
	}
}

func main() {
	arr := []string{"A", "B", "C", "JAYSON", "LouIE"}
	input := make(chan string)

	go func() {
		defer close(input)
		for _, a := range arr {
			input <- a
		}
	}()
	letterAppenderPipeline := AppendLetter("amazing")
	_ = letterAppenderPipeline
	for res := range letterAppenderPipeline(appendLetter2(ToLower(input))) {
		fmt.Println(res)
	}
}
