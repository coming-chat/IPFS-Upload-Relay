package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// Upload file to IPFS
	cid, err := utils.Upload2ForeverLand(reqBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Also upload to W3S
	_, _ = utils.Upload2W3S(bytes.NewReader(reqBytes), "data.json")

	// Return response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"cid":     cid,
		"url":     fmt.Sprintf("ipfs://%s", cid),
		"web2url": utils.AddGateway(cid),
	})

}
