package inits

import (
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"os"
)

func ForeverLand() error {
	var exist bool
	global.ForeverLand_Bucket, exist = os.LookupEnv("4EVERLAND_BUCKET")
	if !exist {
		return fmt.Errorf("env virable 4EVERLAND_BUCKET not found")
	}

	return nil
}
