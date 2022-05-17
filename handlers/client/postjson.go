package client

import (
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func PostJson(ctx *gin.Context) {
	// Validate request
	var req map[string]interface{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Reparse request into bytes
	reqBytes, err := json.Marshal(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Save bytes into tmp file
	tmpFileName := utils.RandomString(16) + ".json"
	if err = os.WriteFile(tmpFileName, reqBytes, 0644); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Upload file to IPFS
	cid, err := ipfsUpload(tmpFileName)
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
		"url":     fmt.Sprintf("ipfs://%s/%s", cid, tmpFileName),
		"web2url": fmt.Sprintf("https://%s.ipfs.dweb.link/%s", cid, tmpFileName),
	})

}
