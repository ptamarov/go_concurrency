package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, totalPizzas int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string // succeeded, failed?
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch

}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1 // delay execution
		fmt.Printf("received order number %d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		totalPizzas++

		fmt.Printf("Making pizza number %d. It will take %d seconds.\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
			success = true
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			success:     success,
			message:     msg,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are trying to make
	var i = 0

	// run forever or until we receive a quit notification (from quit chan)

	// try to make pizzas

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			// The select statement lets a goroutine wait on multiple communication operations.
			// A select blocks until one of its cases can run, then it executes that case.
			// It chooses one at random if multiple are ready.

			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				// unless quit channel is ready to communicate a quit response, select
				// will block or send information to pizzaMaker.data

				// close channels
				close(pizzaMaker.data)
				close(quitChan)
				return // end go routine
			}
		}
		// try to make a single pizza
		// decision structure (made, failed, quit?) with a select statement
	}
}
func main() {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// print a message saying program is starting
	color.Cyan("The pizzeria is open for business.")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background (go routine)
	go pizzeria(pizzaJob)

	// create and run a consumer

	// print out the ending message
}
