package main

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/inits"
	"log"
)

func main() {
	// Init IPFS upstream
	log.Println("Initializing 4everland client...")
	if err := inits.ForeverLand(); err != nil {
		log.Fatalln(err)
	}
	log.Println("4everland initialization complete.")

	// Init redis
	log.Println("Initializing redis client...")
	if err := inits.Redis(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Redis client initialization complete.")

	// Initialize routers
	log.Println("Initializing routers...")
	e := inits.Routers()
	log.Println("Routers initialization complete.")

	// Start server
	log.Println("Starting server...")
	if err := e.Run(); err != nil {
		log.Fatalln(err)
	}
}
