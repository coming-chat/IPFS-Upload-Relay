package inits

import (
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"os"
)

func AwsS3() error {
	var exist bool
	global.AwsS3_Bucket, exist = os.LookupEnv("AWSS3_BUCKET")
	if !exist {
		return fmt.Errorf("env virable AWSS3_BUCKET not found")
	}

	return nil
}
