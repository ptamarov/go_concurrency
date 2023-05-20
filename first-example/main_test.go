package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printSomething(t *testing.T) {
	stdout := os.Stdout

	r, w, _ := os.Pipe() // third arg is an error

	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)
	go printSomething("Testing", &wg)
	wg.Wait()
	_ = w.Close() // saved what was being written in w (Stout), wrote to r

	result, _ := io.ReadAll(r)
	output := string(result) // convert slice of bytes to string

	os.Stdout = stdout // roll back to usual Stdout

	if !strings.Contains(output, "Testing") {
		t.Errorf("Epected to find \"Testing\" but it was not there.")
	}
}
