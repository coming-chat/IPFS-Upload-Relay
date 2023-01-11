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
	global.IPFS_UPLOAD_URL, exist = os.LookupEnv("IPFS_UPLOAD_URL")
	if !exist {
		return fmt.Errorf("env virable IPFS_UPLOAD_URL not found")
	}
	global.PROJECT_ID, exist = os.LookupEnv("PROJECT_ID")
	if !exist {
		return fmt.Errorf("env virable PROJECT_ID not found")
	}
	global.PROJECT_SECRET, exist = os.LookupEnv("PROJECT_SECRET")
	if !exist {
		return fmt.Errorf("env virable PROJECT_SECRET not found")
	}
	global.IPFS_GATEWAY, exist = os.LookupEnv("IPFS_GATEWAY")
	if !exist {
		return fmt.Errorf("env virable IPFS_GATEWAY not found")
	}

	return nil
}
