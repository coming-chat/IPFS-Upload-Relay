package inits

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"os"
	"strings"
)

func W3SClient() error {
	global.W3SAPIKeys = strings.Split(os.Getenv("W3S_TOKEN"), ",")
	return nil
}
