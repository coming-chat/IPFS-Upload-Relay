package main

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/inits"
	"log"
)

func main() {
	// Init
	log.Println("Initializing Web3Storage Client...")
	if err := inits.W3SClient(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Web3Storage Client initialization complete.")

	log.Println("Initializing 4everland client...")
	if err := inits.ForeverLand(); err != nil {
		log.Fatalln(err)
	}
	log.Println("4everland initialization complete.")

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
