package inits

import (
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"os"
)

func Ipfs() error {
	var exist bool
	global.IPFS_URL, exist = os.LookupEnv("IPFS_URL")
	if !exist {
		return fmt.Errorf("env virable IPFS_URL not found")
	}

	return nil
}
