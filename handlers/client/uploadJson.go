package client

import (
	"bytes"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func UploadJson(ctx *gin.Context) {
	log.Print("New binary (json) upload request received")
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to receive request bytes with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	log.Print("Uploading file to IPFS...")
	// UploadFile file to IPFS
	cid, fSize, err := utils.UploadToAwsS3(bytes.NewReader(b), "")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	log.Printf("File uploaded successfully with cid: %s", cid)
	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"cid":      cid,
		"url":      fmt.Sprintf("ipfs://%s", cid),
		"web2url":  utils.AddGateway(cid),
		"fileSize": fSize,
	})

}
