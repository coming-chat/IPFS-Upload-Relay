package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

/******************** Settings ********************/

const (
	VideoDownloadProcessingRecordExpires = 1 * time.Hour
	VideoDownloadFinishedRecordExpires   = 24 * time.Hour
	VideoDownloadStatusRedisKeyTemplate  = "iur:video:%x"
	VideoDownloadTempFilenameTemplate    = "tmp-video-%x"
)

/**************************************************/

type StatusCode string

const (
	VIDEO_STATUS_UPLOAD_SUCCEED    StatusCode = "succeed"
	VIDEO_STATUS_FAILED                       = "failed"
	VIDEO_STATUS_DOWNLOADING                  = "downloading"
	VIDEO_STATUS_UPLOADING_TO_IPFS            = "uploading to IPFS"
)

type VideoUploadStatus struct {
	Status      StatusCode `json:"status"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Message     string     `json:"message"`
	CID         string     `json:"cid"`
	FileSize    int64      `json:"file_size"`
	ContentType string     `json:"content_type"`
}

func UploadVideo(ctx *gin.Context) {

	log.Print("Validating URL ...")
	// POST /video?url=xxxxx
	videoUrl := ctx.Query("url")

	if !utils.VerifyURL(videoUrl) {
		errMsg := "Empty or invalid URL"
		log.Print(errMsg)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": errMsg,
		})
		return
	}

	// Check status in cache
	status, exist, err := getVideoUploadStatus(videoUrl)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get video (%s) status info with error: %s", videoUrl, err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errMsg,
		})
		return
	}

	if !exist {
		// Create new work
		go func() {
			startNewVideoUploadJob(videoUrl)
		}()
		ctx.JSON(http.StatusCreated, gin.H{
			"status": "ok",
			"error":  fmt.Sprintf("New upload work for video (%s) created, please check later.", videoUrl),
		})
	} else {
		// Report work status
		log.Printf("Video (%s) status: %s", videoUrl, status.Status)
		switch status.Status {
		case VIDEO_STATUS_UPLOAD_SUCCEED:
			// Return response
			ctx.JSON(http.StatusOK, gin.H{
				"status":      "ok",
				"cid":         status.CID,
				"url":         fmt.Sprintf("ipfs://%s", status.CID),
				"web2url":     utils.AddGateway(status.CID),
				"updatedAt":   status.UpdatedAt,
				"fileSize":    status.FileSize,
				"contentType": status.ContentType,
			})
			return

		case VIDEO_STATUS_FAILED:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":    "error",
				"message":   status.Message,
				"updatedAt": status.UpdatedAt,
			})
			return

		default:
			ctx.JSON(http.StatusAccepted, gin.H{
				"status":    status.Status,
				"updatedAt": status.UpdatedAt,
			})
			return

		}
	}

}

func startNewVideoUploadJob(videoUrl string) {

	log.Printf("Start downloading video (%s)...", videoUrl)

	// Stage 1: Download

	// youtube-dl -o tmpfile https://...

	setUploadStatus(videoUrl, &VideoUploadStatus{
		Status: VIDEO_STATUS_DOWNLOADING,
	}, VideoDownloadProcessingRecordExpires)

	// Create tmp file
	videoUrlHash := sha256.Sum256([]byte(videoUrl))
	filename := fmt.Sprintf(VideoDownloadTempFilenameTemplate, videoUrlHash)

	// Remove tmp file
	defer os.Remove(filename)

	// Download video
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	dlCmd := exec.Command("./youtube-dl", "-o", filename, videoUrl)
	dlCmd.Stdout = &outBuf
	dlCmd.Stderr = &errBuf
	if err := dlCmd.Run(); err != nil {
		log.Printf("Failed to download video %s with error: %v", videoUrl, err)
		log.Printf(errBuf.String())
		setUploadStatus(videoUrl, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: err.Error(),
		}, VideoDownloadFinishedRecordExpires)
		return
	} else {
		log.Printf("Video %s download successfully", videoUrl)
		log.Printf(outBuf.String())
	}

	// Stage 2: Upload to IPFS

	setUploadStatus(videoUrl, &VideoUploadStatus{
		Status: VIDEO_STATUS_UPLOADING_TO_IPFS,
	}, VideoDownloadProcessingRecordExpires)

	// Reopen to read
	tmpFileR, err := os.Open(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open tmp file with error: %s", err.Error())
		log.Print(errMsg)
		setUploadStatus(videoUrl, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, VideoDownloadFinishedRecordExpires)
		return
	}

	log.Print("Uploading file to IPFS...")
	// UploadFile file to IPFS
	cid, fSize, cType, err := utils.Upload2ForeverLand(tmpFileR)
	_ = tmpFileR.Close() // Close opened file
	if err != nil {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS with error: %s", err.Error())
		log.Print(errMsg)
		setUploadStatus(videoUrl, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, VideoDownloadFinishedRecordExpires)
		return
	}

	// Save result

	log.Printf("File uploaded successfully with cid: %s", cid)
	setUploadStatus(videoUrl, &VideoUploadStatus{
		Status:      VIDEO_STATUS_UPLOAD_SUCCEED,
		CID:         cid,
		FileSize:    fSize,
		ContentType: cType,
	}, VideoDownloadFinishedRecordExpires)
}

func buildCacheKey(videoUrl string) string {
	// Use hash to mark video url
	videoUrlHash := sha256.Sum256([]byte(videoUrl))
	return fmt.Sprintf(VideoDownloadStatusRedisKeyTemplate, videoUrlHash)
}

func setUploadStatus(videoUrl string, status *VideoUploadStatus, expires time.Duration) {
	cacheKey := buildCacheKey(videoUrl)

	status.UpdatedAt = time.Now()

	statusBytes, err := json.Marshal(status)
	if err != nil {
		// ???
		return
	}

	global.Redis.Set(context.Background(), cacheKey, statusBytes, expires)
}

func getVideoUploadStatus(videoID string) (*VideoUploadStatus, bool, error) {
	cacheKey := buildCacheKey(videoID)

	// Check if exist
	exist, err := global.Redis.Exists(context.Background(), cacheKey).Result()
	if err != nil {
		return nil, false, err
	} else if exist == 0 {
		return nil, false, nil // Just not exist
	}

	// Get value
	statusBytes, err := global.Redis.Get(context.Background(), cacheKey).Bytes()
	if err != nil {
		return nil, false, err
	}

	// Parse
	var status VideoUploadStatus
	if err = json.Unmarshal(statusBytes, &status); err != nil {
		return nil, true, err // Exist, but cannot be unmarshalled
	}

	return &status, true, nil

}
