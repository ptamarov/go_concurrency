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
	data chan PizzaOrder // channel to transmit pizza order
	quit chan chan error // channel to STOP production, out of band communication
}

type PizzaOrder struct {
	pizzaNumber int
	message     string // succeeded, failed?
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch // quit channel is now ready to send data
	return <-ch

}

func makePizza(pizzaNumber int) *PizzaOrder {
	// does not interact with other functions
	// all variables are internal

	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(3) + 1 // delay execution
		fmt.Printf("--> Received order number %d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		totalPizzas++

		fmt.Printf("Making  pizza number %d. It will take %d seconds.\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			msg = fmt.Sprintf("<-- Pizza order #%d is ready!", pizzaNumber)
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
	var i = 0 // keep track of which pizza we are trying to make

	// run forever until we receive a quit notification (from quit chan)
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			// The select statement lets a goroutine wait on multiple communication operations.
			// A select blocks until one of its cases can run, then it executes that case.
			// It chooses one at random if multiple are ready.

			case pizzaMaker.data <- *currentPizza:
				// feeds data to pizzaMaker.data channel
				// unless pizzaMaker.quit is ready to communicate
				// a quit response

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
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// print a message saying program is starting
	color.Cyan("The pizzeria is open for business.")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaMaker := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background (go routine)
	go pizzeria(pizzaMaker)

	// create and run a consumer that listens to information
	// coming **into** the pizzMaker.data channel
	// order.pizzaNumber is controlled by variable in line 81

	for order := range pizzaMaker.data {
		//
		if order.pizzaNumber <= NumberOfPizzas {
			if order.success {
				color.Green(order.message)
				color.Green("Order %d is out for delivery!", order.pizzaNumber)
			} else {
				color.Red(order.message)
				color.Red("The customer filed a complaint.")
			}
		} else {
			color.Cyan("Done making pizzas.")
			err := pizzaMaker.Close()
			if err != nil {
				color.Red("Error closing channel.")
			}
		}
	}
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d pizzas, failed to make %d, with %d total attempts.", pizzasMade, pizzasFailed, totalPizzas)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed > 5:
		color.Red("It was not a good day...")
	case pizzasFailed > 3:
		color.Yellow("It was an okay day...")
	case pizzasFailed > 1:
		color.Yellow("It was a pretty good day...")
	default:
		color.Green("It was a great day.")
	}
}
