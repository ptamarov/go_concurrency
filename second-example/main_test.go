package main

import "testing"

func Test_updateMessage(t *testing.T) {
	msg = "Hello, John."

	wg.Add(2)
	go updateMessage("Bye, Bob.")
	go updateMessage("Bye, John.")
	wg.Wait()

	if msg != "Bye, John." {
		t.Error("Incorrect value for msg.")
	}

}
