package main

import (
	"strings"
	"sync"
)

var (
	initialString string
	finalString string
	stringLength int
)

func init() {
	initialString = "Four score and seven years ago our fathers brought forth on this continent, a new nation, " +
		"conceived in Liberty, and dedicated to the proposition " +
		"that all men are created equal."

}

func capitalize(letterChan <-chan string, wg *sync.WaitGroup) <-chan string {
	defer wg.Done()
	outChan := make(chan string)
	go func() {
		for letter := range letterChan {
		outChan <- strings.ToUpper(letter)
	}

	}()
	return outChan
}

func addToFinalStack(letterChan <-chan string, wg sync.WaitGroup) {
	defer wg.Done()
	for letter := range letterChan {
		finalString += letter
	}
}

func main() {
	var wg sync.WaitGroup
	initialBytes :=  []byte(initialString)
	stringLength = len(initialBytes)

	letterChan := make(chan string)



}