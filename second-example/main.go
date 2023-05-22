package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {
	msg = "Berlin."

	wg.Add(2)
	go updateMessage("Prague.")
	go updateMessage("Tokyo.")
	wg.Wait()

	fmt.Println(msg)
}
