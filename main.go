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

	// Initialize routers
	log.Println("Initializing routers...")
	r := inits.Routers()
	log.Println("Routers initialization complete.")

	// Start server
	log.Println("Starting server...")
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
