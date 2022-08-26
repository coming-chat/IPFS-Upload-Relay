package client

import (
	"bytes"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// Upload : Proxy upload request to IPFS
func Upload(ctx *gin.Context) {
	// Receive file from context
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	f, err := file.Open()
	defer f.Close()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Upload file to IPFS
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, f); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	cid, err := utils.Upload2ForeverLand(fileBuffer.Bytes())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Also upload to W3S
	_, _ = f.Seek(0, io.SeekStart)
	_, _ = utils.Upload2W3S(f, file.Filename)

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"cid":     cid,
		"url":     fmt.Sprintf("ipfs://%s", cid),
		"web2url": utils.AddGateway(cid),
	})

}
