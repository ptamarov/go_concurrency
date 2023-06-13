package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup

	words := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		"j",
		"k",
		"l",
		"m",
	}

	wg.Add(len(words))

	for i, word := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, word), &wg)
	}

	wg.Wait()
	wg.Add(1)

	printSomething("Last print", &wg)
}
