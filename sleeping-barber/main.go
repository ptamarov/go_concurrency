package main

import (
	"fmt"
	"math/rand"
	"time"
)

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	loadNames()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// channels to communicate between go routines
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	// barber listens to clientChan and picks customers from there
	shop.addBarber("Frank")

	// channels to stop and wake up routines
	shopClosing := make(chan bool)
	closed := make(chan bool)

	// start a timer that will signal the program to stop after timeOpen seconds
	go func() {
		<-time.After(timeOpen) // block for timeOpen mS
		shopClosing <- true
		shop.closeShopForDay() // blocked by BarbersDone
		closed <- true
	}()

	// continuously add clients until shop is closing
	i := 1

	go func() {
		for {
			// stop routine if shop is closing
			walkingTime := r.Int() % (2 * arrivalRate)

			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(walkingTime)):
				shop.addClient(fmt.Sprintf("Client %d", i))
				i++
			}
		}
	}()

	<-closed // block until the time goroutine finishes
}
