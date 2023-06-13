package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_updateMessage(t *testing.T) {

	wg.Add(1)
	go updateMessage("Testing", &wg)
	wg.Wait()

	if msg != "Testing" {
		t.Error("Expecting to find \"Testing\", but it is not there.")
	}

}

func Test_printMessage(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printMessage("Testing")

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout // roll back to usual Stdout

	if !strings.Contains(output, "Testing") {
		t.Errorf("Expected to find \"Testing\" but it was not there.")
	}
}

func Test_main(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdout // roll back to usual Stdout
	msgs := []string{"Hello, universe!", "Hello, cosmos!", "Hello, world!"}

	for _, m := range msgs {
		if !strings.Contains(output, m) {
			t.Errorf("Epected to find \"Testing\" but it was not there.")
		}
	}

}
