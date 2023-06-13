package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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
