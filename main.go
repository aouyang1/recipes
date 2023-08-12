package main

import (
	"log"
)

func main() {
	cf, err := NewCalendarFetcher()
	if err != nil {
		log.Fatal(err)
	}
	if err := cf.Fetch(); err != nil {
		log.Fatal(err)
	}
}
