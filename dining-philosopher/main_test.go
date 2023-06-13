package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		orderFinished = []string{}
		dine()
		if len(orderFinished) != 5 {
			t.Errorf("Incorrect length of slice: expected 5 but got %d", len(orderFinished))
		}
	}
}

func Test_dineWithDelays(t *testing.T) {
	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"Zero delay", 0 * time.Second},
		{"Quarter second", 250 * time.Millisecond},
		{"Half second", 500 * time.Millisecond},
	}

	for _, test := range theTests {
		orderFinished = []string{}

		eatTime = test.delay
		thinkTime = test.delay

		dine()
		if len(orderFinished) != 5 {
			t.Errorf("%s: Incorrect length of slice: expected 5 but got %d", test.name, len(orderFinished))
		}
	}
}
