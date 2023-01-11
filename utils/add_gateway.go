package utils

import (
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
)

func AddGateway(cid string) string {
	return fmt.Sprintf(global.IPFS_GATEWAY, cid)
}
