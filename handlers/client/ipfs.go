package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
)

type Web3StorageResponse struct {
	Cid     string `json:"cid"`
	Name    string `json:"name"`    // Only on failure
	Message string `json:"message"` // Only on failure
}

func ipfsUpload(r io.Reader, filename string) (string, error) {
	// Prepare
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	// Copy file to Writer
	_, err = io.Copy(fileWriter, r)
	if err != nil {
		log.Println("Error copying file to Writer:", err)
		return "", err
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

	if w3sResponse.Cid == "" {
		return "", fmt.Errorf(w3sResponse.Message)
	}

	// Return response
	return w3sResponse.Cid, nil
}
