package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Record struct {
	sync.Mutex // embed

	buf string
	cond *sync.Cond

	writers []io.Writer
}

func NewRecord(writers ...io.Writer) *Record {
	r := &Record{ writers: writers }
	r.cond = sync.NewCond(r)
	return r
}

func (r *Record) Prompt() {
	for {
		fmt.Printf(":> ")
		var s string
		fmt.Scanf("%s", &s)

		r.Lock()
		r.buf = s
		r.Unlock()

		r.cond.Broadcast()
	}
}

func (r *Record) Start() error {
	f := func(w io.Writer) {
		for {
			r.Lock()
			r.cond.Wait() // waiting for the signal/broadcast
			 		      // each time there is a new data, all those waiting goroutine get notified via Broadcast()
			fmt.Fprintf(w, "%s\n", r.buf)
			r.Unlock()
		}
	}
	for i := range r.writers {
		go f(r.writers[i]) // go for background
	}
	return nil
}

func main() {
	f, err := os.Create("cond.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f2, err := os.Create("cond2.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()
	r := NewRecord(f, f2, os.Stdout)
	r.Start()
	r.Prompt()
}

