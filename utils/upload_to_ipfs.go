package utils

//
//import (
//	"context"
//	shell "github.com/ipfs/go-ipfs-api"
//	files "github.com/ipfs/go-ipfs-files"
//	"io"
//	"os"
//	"strings"
//)
//
//func UploadToIpfs(r io.ReadSeeker) (string, int64, error) {
//	ipfsUploadUrl, exist := os.LookupEnv("IPFS_UPLOAD_URL")
//	if !exist {
//		panic("IPFS_UPLOAD_URL is not exist")
//	}
//	ipfs := shell.NewShell(ipfsUploadUrl)
//	fr := files.NewReaderFile(r)
//	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
//	fileReader := files.NewMultiFileReader(slf, true)
//	ipfs.Add()
//	fileData := make(map[string]interface{})
//	rb := ipfs.Request("add")
//	err := rb.Body(fileReader).Exec(context.Background(), &fileData)
//	if err != nil {
//		return "", 0, err
//	}
//
//	err = ipfs.FilesCp(context.Background(), "/ipfs/"+fileData["Hash"].(string), "/"+fileData["Name"].(string))
//	if err != nil {
//		return "", 0, err
//	}
//
//	return strings.ReplaceAll(fileData["Hash"].(string), "\"", ""), fileData["Size"].(int64), nil
//}
