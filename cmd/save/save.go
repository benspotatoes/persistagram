package main

import (
	"log"
	"os"

	"github.com/benspotatoes/persistagram/backend"
)

func main() {
	backend := backend.NewBackend()

	if len(os.Args) != 2 {
		log.Fatal("Invalid arguments")
	}
	backend.Save(os.Args[1])
}
