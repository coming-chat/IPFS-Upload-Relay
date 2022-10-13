package client

import (
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	y2bCli = youtube.Client{}
)

func UploadVideo(ctx *gin.Context) {
	// Use link to upload
	log.Print("New video upload request received")

	log.Print("Validating URL ...")
	// POST /video?url=xxxxx
	videoUrl := ctx.Query("url")

	if !utils.VerifyURL(videoUrl) {
		errMsg := "Empty or invalid URL"
		log.Print(errMsg)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	log.Print("Start downloading video...")
	// Download
	video, err := y2bCli.GetVideo(videoUrl)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get video with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := y2bCli.GetStream(video, &formats[0])
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create video stream with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("Failed to create video stream with error: %s", err.Error()),
		})
		return
	}

	// Create tmp file
	filename := fmt.Sprintf("tmp-video-%d", time.Now().UnixNano())
	tmpFileW, err := os.Create(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create tmp file with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	defer os.Remove(filename) // Delete tmp file

	// Download video
	_, err = io.Copy(tmpFileW, stream)
	_ = tmpFileW.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to download video with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	// Reopen to read
	tmpFileR, err := os.Open(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open tmp file with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  errMsg,
		})
		return
	}

	log.Print("Uploading file to IPFS...")
	// UploadFile file to IPFS
	cid, err := utils.Upload2ForeverLand(tmpFileR)
	_ = tmpFileR.Close() // Close opened file
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
		"status":  "ok",
		"cid":     cid,
		"url":     fmt.Sprintf("ipfs://%s", cid),
		"web2url": utils.AddGateway(cid),
	})

}
