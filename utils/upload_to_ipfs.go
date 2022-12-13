package utils

import (
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"os"
	"strings"
)

func UploadToIpfs(r io.ReadSeeker) (string, int64, error) {
	ipfsUploadUrl, exist := os.LookupEnv("IPFS_UPLOAD_URL")
	if !exist {
		panic("IPFS_UPLOAD_URL is not exist")
	}
	ipfs := shell.NewShell(ipfsUploadUrl)
	add, err := ipfs.Add(r, shell.OnlyHash(true))
	if err != nil {
		return "", 0, err
	}
	return strings.ReplaceAll(add, "\"", ""), 0, nil
}
