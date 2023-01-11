package utils

import (
	"bytes"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	ipfs "github.com/ipfs/go-ipfs-api"
)

func GetIPFSCid(file []byte) (string, error) {
	shell := ipfs.NewShell(global.IPFS_URL)
	cid, err := shell.Add(bytes.NewReader(file), ipfs.CidVersion(1), ipfs.OnlyHash(true))
	if err != nil {
		return "", err
	}

	return cid, nil
}
