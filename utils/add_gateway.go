package utils

import "fmt"

func AddGateway(cid string) string {
	return fmt.Sprintf("https://4everland.io/ipfs/%s", cid)
}
