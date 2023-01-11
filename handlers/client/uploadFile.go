package client

import (
	"bytes"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// UploadFile : Proxy upload request to IPFS
func UploadFile(ctx *gin.Context) {
	// Receive file from context
	log.Print("New file upload request received")

	timeUploadStart := time.Now()

	file, err := ctx.FormFile("file")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get input form file with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	timeUpload2Relay := time.Now()

	log.Print("Opening form file...")
	f, err := file.Open()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open input form file with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}
	defer f.Close()

	log.Print("Uploading file to IPFS...")
	// UploadFile file to IPFS
	var (
		wg           = sync.WaitGroup{}
		err1, err2   error
		cidIpfs, cid string
		fSize        int64
	)
	wg.Add(2)
	fileBuf, err := io.ReadAll(f)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}
	_, _ = f.Seek(0, io.SeekStart)
	go func() {
		cidIpfs, _, err1 = utils.UploadToIpfs(bytes.NewReader(fileBuf))
		wg.Done()
	}()
	go func() {
		cid, fSize, err2 = utils.UploadToAwsS3(f)
		wg.Done()
	}()

	wg.Wait()
	if err1 != nil || err2 != nil {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS")
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	if cidIpfs != cid {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	timeUpload2IPFS := time.Now()

	log.Printf("File uploaded successfully with cid: %s", cid)
	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"cid":      cid,
		"url":      fmt.Sprintf("ipfs://%s", cid),
		"web2url":  utils.AddGateway(cid),
		"fileSize": fSize,
	})

	log.Printf(
		"Time consumption: %dms to relay, %dms to IPFS",
		timeUpload2Relay.Sub(timeUploadStart).Microseconds(),
		timeUpload2IPFS.Sub(timeUpload2Relay).Microseconds(),
	)

}
