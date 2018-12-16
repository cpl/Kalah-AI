package main

import (
	"log"
)

func main() {
	go startGigel()

	startPlayer()
}

func startPlayer() {
	// create player
	p, err := newPlayer("127.0.0.1", "12340")
	if err != nil {
		log.Fatalln(err)
	}

	p.start()
}

func startGigel() {
	// use standard default weights
	weights := [7]int{
		15,
		15,
		35,
		100,
		40,
		-80,
		-20,
	}

	// create gigel on 127.0.0.1
	g, err := newGigel("127.0.0.1", "12345", 4, weights)
	if err != nil {
		log.Fatalln(err)
	}

	// start gigel
	_, err = g.start()
	if err != nil {
		log.Fatalln(err)
	}
}
