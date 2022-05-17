package inits

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/web3-storage/go-w3s-client"
	"os"
)

func W3SClient() error {
	var err error
	global.W3SClient, err = w3s.NewClient(w3s.WithToken(os.Getenv("W3S_TOKEN")))
	return err
}
