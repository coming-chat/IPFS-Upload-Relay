package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

/******************** Settings ********************/

const (
	VideoDownloadProcessingRecordExpires = 1 * time.Hour
	VideoDownloadFinishedRecordExpires   = 24 * time.Hour
	VideoDownloadStatusRedisKeyTemplate  = "iur:video:%s"
	VideoDownloadTempFilenameTemplate    = "tmp-video-%d"
)

/**************************************************/

type StatusCode int

const (
	VIDEO_STATUS_UPLOAD_SUCCEED    StatusCode = 1
	VIDEO_STATUS_FAILED                       = -1
	VIDEO_STATUS_LOADING_METADATA             = 2
	VIDEO_STATUS_DOWNLOADING                  = 3
	VIDEO_STATUS_UPLOADING_TO_IPFS            = 4
)

type VideoUploadStatus struct {
	Status    StatusCode `json:"status"`
	UpdatedAt time.Time  `json:"updated_at"`
	Message   string     `json:"message"`
	CID       string     `json:"cid"`
}

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
			"status":  "error",
			"message": errMsg,
		})
		return
	}

	log.Print("Extracting video ID...")
	// Download
	videoID := videoUrl

	// Check status in cache
	status, exist, err := getVideoUploadStatus(videoID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get video (%s) status info with error: %s", videoID, err.Error())
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
			startNewVideoUploadJob(videoID)
		}()
		ctx.JSON(http.StatusCreated, gin.H{
			"status": "ok",
			"error":  fmt.Sprintf("New upload work for video (%s) created, please check later.", videoID),
		})
	} else {
		// Report work status
		var videoStatusStr string
		switch status.Status {
		case VIDEO_STATUS_UPLOAD_SUCCEED:
			// Return response
			ctx.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"cid":       status.CID,
				"url":       fmt.Sprintf("ipfs://%s", status.CID),
				"web2url":   utils.AddGateway(status.CID),
				"updatedAt": status.UpdatedAt,
			})
			return

		case VIDEO_STATUS_FAILED:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":    "error",
				"message":   status.Message,
				"updatedAt": status.UpdatedAt,
			})
			return

		case VIDEO_STATUS_LOADING_METADATA:
			videoStatusStr = "loading metadata"
		case VIDEO_STATUS_DOWNLOADING:
			videoStatusStr = "downloading"
		case VIDEO_STATUS_UPLOADING_TO_IPFS:
			videoStatusStr = "uploading to IPFS"

		default:
			// ???
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":    "unknown",
				"updatedAt": status.UpdatedAt,
			})
			return

		}

		log.Printf("Video (%s) status: %s", videoID, videoStatusStr)
		ctx.JSON(http.StatusAccepted, gin.H{
			"status":    videoStatusStr,
			"updatedAt": status.UpdatedAt,
		})
	}

}

func startNewVideoUploadJob(videoID string) {

	log.Printf("Start downloading video (%s)...", videoID)

	// Stage 2: Download from Y2B

	// youtube-dl -o tmpfile https://...

	setY2BUploadStatus(videoID, &VideoUploadStatus{
		Status: VIDEO_STATUS_DOWNLOADING,
	}, VideoDownloadProcessingRecordExpires)

	// Create tmp file
	filename := fmt.Sprintf(VideoDownloadTempFilenameTemplate, time.Now().UnixMicro())
	tmpFileW, err := os.Create(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create tmp file with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, VideoDownloadFinishedRecordExpires)
		return
	}

	defer os.Remove(filename) // Delete tmp file

	// Download video
	_, err = io.Copy(tmpFileW, stream)
	_ = tmpFileW.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to download video with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, VideoDownloadFinishedRecordExpires)
		return
	}

	// Stage 3: Upload to IPFS

	setY2BUploadStatus(videoID, &VideoUploadStatus{
		Status: VIDEO_STATUS_UPLOADING_TO_IPFS,
	}, VideoDownloadProcessingRecordExpires)

	// Reopen to read
	tmpFileR, err := os.Open(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open tmp file with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, VideoDownloadFinishedRecordExpires)
		return
	}

	log.Print("Uploading file to IPFS...")
	// UploadFile file to IPFS
	cid, err := utils.Upload2ForeverLand(tmpFileR)
	_ = tmpFileR.Close() // Close opened file
	if err != nil {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &VideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, VideoDownloadFinishedRecordExpires)
		return
	}

	// Save result

	log.Printf("File uploaded successfully with cid: %s", cid)
	setY2BUploadStatus(videoID, &VideoUploadStatus{
		Status: VIDEO_STATUS_UPLOAD_SUCCEED,
		CID:    cid,
	}, VideoDownloadFinishedRecordExpires)
}

func buildCacheKey(videoID string) string {
	return fmt.Sprintf(VideoDownloadStatusRedisKeyTemplate, videoID)
}

func setY2BUploadStatus(videoID string, status *VideoUploadStatus, expires time.Duration) {
	cacheKey := buildCacheKey(videoID)

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
