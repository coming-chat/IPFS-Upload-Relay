package client

import (
	"context"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"log"
	"os"
)

func ipfsUpload(filename string) (string, error) {
	// Reopen file for IPFS upload
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	// Proxy to IPFS
	cid, _ := global.W3SClient.Put(context.Background(), f)
	log.Println("File upload success", cid)

	//// Close tmp file // Web3Storage SDK will close the file after upload, so we don't need to close it here
	//if err = f.Close(); err != nil {
	//	log.Println("Unable to close tmp file", err.Error())
	//}

	// Unlink tmp file
	if err = os.Remove(filename); err != nil {
		log.Println("Unable to remove tmp file", err.Error())
	}

	// Return response
	return cid.String(), nil
}
