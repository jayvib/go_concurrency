package main

import (
	"log"
	"fmt"
	"os"
	"github.com/go_concurrency/error-propagation/intermediate"
	"github.com/go_concurrency/error-propagation/error-prop"
)

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	errfile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(error_prop.WrapError(err, "somethings wrong when opening the file %s", "logs.log"))
	}


	log.SetOutput(errfile)
	log.SetFlags(log.Ltime|log.LUTC)

	err = intermediate.RunJob("1")
	if err != nil {
		msg := "There was an unexpected issue: please report this as a bug."
		if _, ok := err.(intermediate.IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg) // unexpected error
	}
}
