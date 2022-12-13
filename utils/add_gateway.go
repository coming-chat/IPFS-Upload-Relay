package utils

import (
	"fmt"
	"os"
)

func AddGateway(cid string) string {
	ipfsUrl, exist := os.LookupEnv("IPFS_URL")
	if !exist {
		panic("IPFS_URL is not exist")
	}
	return fmt.Sprintf(ipfsUrl, cid)
}
