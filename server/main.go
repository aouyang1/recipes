package main

import (
	"log"
)

func main() {
	s, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
