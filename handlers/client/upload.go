package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	cid, err := ipfsUpload(f, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"cid":     cid,
		"url":     fmt.Sprintf("ipfs://%s", cid),
		"web2url": fmt.Sprintf("https://%s.ipfs.cf-ipfs.com", cid),
	})

}
