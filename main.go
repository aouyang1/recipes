package main

import (
	"log"
)

const (
	calendarSummary = "Dinner Plans"
)

func main() {
	cf, err := NewCalendarFetcher(calendarSummary)
	if err != nil {
		log.Fatal(err)
	}
	if err := cf.Fetch(); err != nil {
		log.Fatal(err)
	}
}
