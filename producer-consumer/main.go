package main

const numberOfPizzas = 10

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

func main() {
	// seed random number generator

	// print a message saying program is starting

	// create a producer

	// run the producer in the background (go routine)

	// create and run a consumer

	// print out the ending message
}
