package main

import (
	"fmt"
	"sync"
	"time"
)

// Philosopher stores information about a philosopher.
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// philosophers is a list of all philosophers.
var philosophers = []Philosopher{
	{name: "Plato", rightFork: 4, leftFork: 0},
	{name: "Socrates", rightFork: 0, leftFork: 1},
	{name: "Aristotle", rightFork: 1, leftFork: 2},
	{name: "Pascal", rightFork: 2, leftFork: 3},
	{name: "Locke", rightFork: 3, leftFork: 4},
}

// define some variables.

var hunger = 3 // how many times does a person eat?
var eatTime = 1 * time.Second
var thinkTime = 1 * time.Second

// *** for challenge
var outQueue = []string{}
var orderOut sync.Mutex

func main() {
	// print a welcome message.
	fmt.Println("Dining Philosophers Problem.")
	fmt.Println("----------------------------")

	fmt.Println("The table is empty.")

	// start the meal.
	dine()

	// print out finished message.
	fmt.Println("The table is empty.")
}

func dine() {
	// wait group to wait for everyone to be done eating.
	mealwg := &sync.WaitGroup{}
	mealwg.Add(len(philosophers))

	// wait group to wait for everyone to be seated.
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all five forks.
	var forks = make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal.
	for i := 0; i < len(philosophers); i++ {
		// fire off a go routine for the current philosopher

		input := diningInput{
			philosopher: philosophers[i],
			forks:       forks,
			meal:        mealwg,
			seated:      seated}

		go diningProblem(input)
	}

	mealwg.Wait()

	fmt.Println(outQueue)
}

type diningInput struct {
	philosopher Philosopher
	forks       map[int]*sync.Mutex
	meal        *sync.WaitGroup
	seated      *sync.WaitGroup
}

func recordExit(p Philosopher, exit *sync.Mutex) {
	exit.Lock()
	outQueue = append(outQueue, p.name)
	exit.Unlock()
}

func diningProblem(d diningInput) {
	defer d.meal.Done()

	// seat the philosopher at the table
	fmt.Printf("%s is seated at the table.\n", d.philosopher.name)

	d.seated.Done() // will decrease seated by 1 as each philospher sits down
	d.seated.Wait() // wait for all philosophers to sit

	// eat three times
	for i := hunger; i > 0; i-- {

		if d.philosopher.leftFork > d.philosopher.rightFork { // prevent all philosophers from taking all left forks at once
			d.forks[d.philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", d.philosopher.name)
			d.forks[d.philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", d.philosopher.name)
		} else {
			d.forks[d.philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", d.philosopher.name)
			d.forks[d.philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", d.philosopher.name)
		}
		// get a lock on both forks

		// now philosopher has both forks
		fmt.Printf("\t%s has both forks and is eating.\n", d.philosopher.name)
		time.Sleep(eatTime)

		// philosopher is thiking
		fmt.Printf("\t%s is thinking.\n", d.philosopher.name)
		time.Sleep(thinkTime)

		d.forks[d.philosopher.leftFork].Unlock()
		d.forks[d.philosopher.rightFork].Unlock()

		fmt.Printf("\t%s puts down the forks.\n", d.philosopher.name)
	}

	fmt.Printf("%s is satisfied.\n", d.philosopher.name)
	fmt.Printf("--> %s left the table.\n", d.philosopher.name)
	recordExit(d.philosopher, &orderOut)
}
