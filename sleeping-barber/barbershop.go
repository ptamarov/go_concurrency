package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Barber string

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

// addBarber reads from ClientsChan continuously until it is closed.
func (shop *BarberShop) addBarber(b Barber) {
	shop.NumberOfBarbers++

	go func() {
		for {
			if len(shop.ClientsChan) == 0 {
				fmt.Printf("%s says: no customers for me.\n", b)
			}

			d, customersLeft := <-shop.ClientsChan

			if customersLeft {
				shop.CutHair(b, d)
			} else {
				// ClientsChan only closed if shop is closed. Go home!
				// Unblocks call for closeShopForDay(),
				// but only after customersLeft returns false
				shop.sendBarberHome(b)
				return
			}
		}
	}()
}

func (shop *BarberShop) CutHair(barber Barber, client string) {
	color.Green("%s is cutting %s's hair", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is done cutting %s's hair", barber, client)
}

func (shop *BarberShop) sendBarberHome(b Barber) {
	color.Green("%s is going home", b)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day")
	close(shop.ClientsChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan // will block until all barbers have sent true to this channel
	}

	close(shop.BarbersDoneChan)
	color.Green("The barber shop is now closed for the day. Everyone has gone home.")
	color.Green("------------------------------------------------------------------")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("*** %s arrives ***", client)

	if shop.Open {
		select {
		case shop.ClientsChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default:
			color.Red("The waiting room is full, so %s leaves.", client)

		}
	} else {
		color.Red("The shop is already closed, so %s leaves.")
	}
}
