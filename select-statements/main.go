package main

import (
	"fmt"
	"time"
)

func server1(ch chan string) {
	for {
		time.Sleep(2 * time.Second)
		ch <- "This is from Server 1."
	}
}

func server2(ch chan string) {
	for {
		time.Sleep(1 * time.Second)
		ch <- "This is from Server 2."
	}
}

func main() {
	fmt.Println("Select statements.")
	fmt.Println("------------------")

	channel1 := make(chan string)
	channel2 := make(chan string)

	go server1(channel1)
	go server2(channel2)

	// the point of the case list below is that
	// select chooses cases at random from the list
	// of cases that are ready!

	for {
		select {
		case s1 := <-channel1:
			fmt.Println("Case one:", s1)
		case s2 := <-channel1:
			fmt.Println("Case two:", s2)
		case s3 := <-channel2:
			fmt.Println("Case three:", s3)
		case s4 := <-channel2:
			fmt.Println("Case three:", s4)
		}
	}

}
