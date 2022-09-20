package client

import (
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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	defer f.Close()

	// Upload file to IPFS
	cid, err := utils.Upload2ForeverLand(f)
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
