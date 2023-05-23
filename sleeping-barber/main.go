package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// first try. code is buggy: works as expected but instead of exiting
// all go routines go to sleep at end, so main exits with error

type Customer struct {
	Name            string
	AppointmentTime int
}

var queueList []string
var waitList sync.Mutex

func main() {
	loadNames()

	// channels to communicate between go routines
	shopDoor := make(chan Customer)
	customerQueue := make(chan Customer, 4)

	// channels to stop and wake up routines
	closeShop := make(chan chan error)
	stopBarber := make(chan chan error)

	openTime := 30

	// fetch customers from shopDoor and attempt to put them in the queue. if the queue
	// is full (or closed), then tell them to go away
	go barberShop(shopDoor, customerQueue, closeShop)

	// start a timer that will signal the program to stop after openTime seconds
	go timeToClose(openTime, closeShop, stopBarber)

	// barber listens to queue and picks customers from there
	go barber(customerQueue, stopBarber)

	// run a customer factory that sends customer to the shopDoor (channel)
	for {
		customerFactory(shopDoor)
	}
}

// customerFactory creates a continuous stream of customers and sends them to the shop

func barber(customerQueue chan Customer, stopBarber chan chan error) {
	for {
		select {

		// while shop is not closed, try to fetch a customer from the queue
		case d := <-customerQueue:

			waitList.Lock()
			queueList = queueList[1:]
			color.Cyan("-----------------BARBER-------------------------")
			fmt.Printf("\tCustomer %s sits in my chair.\n", d.Name)
			fmt.Println("\tI see that", queueList, "are waiting.")
			waitList.Unlock()

			fmt.Printf("\tI am cutting %s's hair. It will take %d seconds.\n", d.Name, d.AppointmentTime)
			time.Sleep(time.Duration(d.AppointmentTime) * time.Second)
			fmt.Printf("\tI am done cutting %s' hair.\n", d.Name)
			color.Cyan("-------------------------------------------------")

			// check if signal was given to close down shop
		case <-stopBarber:
			// process last customers if any
			for d := range customerQueue {
				color.Cyan("----------------BARBER--------------------------")
				fmt.Printf("\t%s was still in the queue. His haircut will take %d seconds.\n", d.Name, d.AppointmentTime)
				time.Sleep(time.Duration(d.AppointmentTime) * time.Second)
				fmt.Printf("\tI am done cutting %s' hair.\n", d.Name)
				color.Cyan("-------------------------------------------------")
			}
			color.Cyan("--------------BARBER GOES HOME. BYE!-------------------------")
			return
		}
	}
}

func customerFactory(d chan Customer) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// get a random number to get a duration
	newDuration := 5 + (r.Intn(2) + 1)

	// get  random numbers to fetch a name and wait
	newIndex := r.Intn(100)
	walkingTime := (r.Intn(2) + 1)

	// create a customer
	customer := Customer{Name: customerNames[newIndex], AppointmentTime: newDuration}

	// give fair warning
	color.Red("----------------------FACTORY--------------------------------")
	fmt.Printf("*** A customer named %s entered the shop... ***\n", customer.Name)
	fmt.Printf("*** They will be there in %d seconds ***\n", walkingTime)
	color.Red("-------------------------------------------------------------")

	// emulate delay when arriving to shop
	time.Sleep(time.Duration(walkingTime) * time.Second)

	d <- customer

}

// queue manager listens continuously to door and tries to put customers in a queue
// or sends them back home if the queue is full
func barberShop(door chan Customer, customerQueue chan Customer, closeShop chan chan error) {
	for {
		d := <-door // read from door, unbuffered channel
		customerName := d.Name

		color.Green("-----------------------BARBER SHOP------------------------------")
		fmt.Printf("Queue manager is trying to sort customer %s.\n", customerName)

		select {
		case <-closeShop:
			color.Green("-----------------------BARBER SHOP------------------------------")
			fmt.Println("BARBER SHOP: Shop is closed. No more customers are accepted.")
			color.Green("----------------------SHOP CLOSES. BYE!-------------------------")
			return

		case customerQueue <- d: // try to send to queue, capacity 4
			waitList.Lock()
			queueList = append(queueList, customerName)
			waitList.Unlock()

			fmt.Printf("BARBER SHOP: Customer %s has joined the queue.\n", customerName)
			fmt.Println("BARBER SHOP: The queue is now", queueList)
			color.Green("---------------------------------------------------------------")

		default:
			fmt.Printf("BARBER SHOP: Sorry %s, the queue is full! Come back later.\n", customerName)
			color.Green("---------------------------------------------------------------")
		}
	}
}

// timeToClose waits openTime seconds and then sends a signal to close shop
func timeToClose(openTime int, closeChannel, stopBarber chan chan error) {
	color.Red("---------------------------TIMER------=------------------------")
	fmt.Printf("\tThe shop will be open for %d seconds.\n", openTime)

	time.Sleep(time.Duration(openTime) * time.Second)
	fmt.Println("Time is up!")
	color.Red("----------------------- TIMER STOPS. BYE!--------------------------")

	signal1 := make(chan error)
	signal2 := make(chan error)
	closeChannel <- signal1
	stopBarber <- signal2

}

var customerNames []string

// loadNames populates customerNames with 100 common English first names.
func loadNames() {
	file, err1 := os.Open("names.txt")

	if err1 != nil {
		log.Fatal("Error reading name file:", err1)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s := fmt.Sprintln(scanner.Text())
		customerNames = append(customerNames, strings.TrimSuffix(s, "\n"))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file.Close()
}
