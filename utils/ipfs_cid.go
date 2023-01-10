package utils

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func GetIPFSCid(file []byte) (string, error) {
	p := cid.Prefix{
		Version:  1,                  // 自定义选择版本 取值0 或者  1
		Codec:    cid.Raw,            // prtobuf
		MhType:   multihash.SHA2_256, // sha2-256
		MhLength: multihash.DefaultLengths[multihash.SHA2_256],
	}

	fileCid, err := p.Sum(file)
	if err != nil {
		return "", err
	}

	return fileCid.String(), nil
}
