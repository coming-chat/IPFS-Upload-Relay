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

	// Save tmp file
	if err = ctx.SaveUploadedFile(file, file.Filename); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Upload file to IPFS
	cid, err := ipfsUpload(file.Filename)
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
		"url":     fmt.Sprintf("ipfs://%s/%s", cid, file.Filename),
		"web2url": fmt.Sprintf("https://%s.ipfs.cf-ipfs.com/%s", cid, file.Filename),
	})

}
