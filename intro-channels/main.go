package main

import (
	"fmt"
	"strings"
)

func shout(ping <-chan string, pong chan<- string) {
	// receive only channel: <-chan
	// send only channel: chan<-
	for {
		s := <-ping                                      // receives from ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s)) // sends to pong
	}
}

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and press ENTER (Q to quit)")

	for {
		// print a prompt
		fmt.Print("-> ")
		var userInput string

		_, _ = fmt.Scanln(&userInput) //

		if strings.ToLower(userInput) == "q" {
			break
		}

		ping <- userInput
		// wait for a reponse from pong
		response := <-pong
		fmt.Printf("Response: %s\n", response)

	}

	fmt.Println("All done, closing channels.")
	close(ping)
	close(pong)
}
