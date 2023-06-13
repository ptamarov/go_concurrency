package main

import (
	"fmt"
	"sync"
)

var msg1 string
var wg1 sync.WaitGroup

func updateMessageWithMutex(s string, m *sync.Mutex) {
	defer wg1.Done()
	m.Lock() // get exclusive access to msg1 (thread safe operation)
	msg1 = s
	m.Unlock()
}

func mainMutex() {
	msg1 = "Berlin."
	var mutex sync.Mutex

	wg1.Add(2)

	go updateMessageWithMutex("Prague.", &mutex)
	go updateMessageWithMutex("Tokyo.", &mutex)

	wg1.Wait()

	fmt.Println(msg1)
}
