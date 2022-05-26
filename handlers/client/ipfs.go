package client

import (
	"bytes"
	"encoding/json"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
)

type Web3StorageResponse struct {
	Cid string `json:"cid"`
}

func ipfsUpload(filename string) (string, error) {
	// Prepare
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	// Reopen file for IPFS upload
	f, err := os.Open(filename)
	//defer f.Close()
	if err != nil {
		log.Println("Error opening file:", err)
		return "", err
	}

	// Copy file to Writer
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		log.Println("Error copying file to Writer:", err)
		return "", err
	}

	// Close and unlink tmp file
	f.Close()
	if err = os.Remove(filename); err != nil {
		log.Println("Unable to remove tmp file", err.Error())
	}

	// Proxy to IPFS
	var w3sResponse Web3StorageResponse
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	req, _ := http.NewRequest("POST", "https://api.web3.storage/upload", bodyBuffer)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+global.W3SAPIKeys[rand.Intn(len(global.W3SAPIKeys))])
	res, _ := (&http.Client{}).Do(req)
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&w3sResponse)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return "", err
	}
	log.Println("File upload success", w3sResponse.Cid)

	// Return response
	return w3sResponse.Cid, nil
}
