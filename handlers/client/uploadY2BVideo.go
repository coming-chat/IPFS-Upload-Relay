package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
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

/******************** Settings ********************/

const (
	Y2BVideoDownloadProcessingRecordExpires = 1 * time.Hour
	Y2BVideoDownloadFinishedRecordExpires   = 24 * time.Hour
	Y2BVideoDownloadStatusRedisKeyTemplate  = "iur:y2b:%s"
	Y2BVideoDownloadTempFilenameTemplate    = "tmp-y2b-%s-%d"
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

type Y2BVideoUploadStatus struct {
	Status    StatusCode `json:"status"`
	UpdatedAt time.Time  `json:"updated_at"`
	Message   string     `json:"message"`
	CID       string     `json:"cid"`
}

func UploadY2BVideo(ctx *gin.Context) {
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
	videoID, err := youtube.ExtractVideoID(videoUrl)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get video ID with error: %s", err.Error())
		log.Print(errMsg)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errMsg,
		})
		return
	}

	// Check status in cache
	status, exist, err := getY2BUploadStatus(videoID)
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
			startNewY2BUploadJob(videoID)
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

func startNewY2BUploadJob(videoID string) {

	log.Printf("Start downloading video (%s)...", videoID)

	// Stage 1: Get video info
	setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
		Status: VIDEO_STATUS_LOADING_METADATA,
	}, Y2BVideoDownloadProcessingRecordExpires)

	video, err := y2bCli.GetVideo(videoID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get video info with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, Y2BVideoDownloadFinishedRecordExpires)
		return
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := y2bCli.GetStream(video, &formats[0])
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create video stream with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, Y2BVideoDownloadFinishedRecordExpires)
		return
	}

	// Stage 2: Download from Y2B
	setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
		Status: VIDEO_STATUS_DOWNLOADING,
	}, Y2BVideoDownloadProcessingRecordExpires)

	// Create tmp file
	filename := fmt.Sprintf(Y2BVideoDownloadTempFilenameTemplate, video.ID, time.Now().UnixMicro())
	tmpFileW, err := os.Create(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create tmp file with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, Y2BVideoDownloadFinishedRecordExpires)
		return
	}

	defer os.Remove(filename) // Delete tmp file

	// Download video
	_, err = io.Copy(tmpFileW, stream)
	_ = tmpFileW.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to download video with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, Y2BVideoDownloadFinishedRecordExpires)
		return
	}

	// Stage 3: Upload to IPFS

	setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
		Status: VIDEO_STATUS_UPLOADING_TO_IPFS,
	}, Y2BVideoDownloadProcessingRecordExpires)

	// Reopen to read
	tmpFileR, err := os.Open(filename)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to open tmp file with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, Y2BVideoDownloadFinishedRecordExpires)
		return
	}

	log.Print("Uploading file to IPFS...")
	// UploadFile file to IPFS
	cid, err := utils.Upload2ForeverLand(tmpFileR)
	_ = tmpFileR.Close() // Close opened file
	if err != nil {
		errMsg := fmt.Sprintf("Failed to upload file to IPFS with error: %s", err.Error())
		log.Print(errMsg)
		setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
			Status:  VIDEO_STATUS_FAILED,
			Message: errMsg,
		}, Y2BVideoDownloadFinishedRecordExpires)
		return
	}

	// Save result

	log.Printf("File uploaded successfully with cid: %s", cid)
	setY2BUploadStatus(videoID, &Y2BVideoUploadStatus{
		Status: VIDEO_STATUS_UPLOAD_SUCCEED,
		CID:    cid,
	}, Y2BVideoDownloadFinishedRecordExpires)
}

func buildCacheKey(videoID string) string {
	return fmt.Sprintf(Y2BVideoDownloadStatusRedisKeyTemplate, videoID)
}

func setY2BUploadStatus(videoID string, status *Y2BVideoUploadStatus, expires time.Duration) {
	cacheKey := buildCacheKey(videoID)

	status.UpdatedAt = time.Now()

	statusBytes, err := json.Marshal(status)
	if err != nil {
		// ???
		return
	}

	global.Redis.Set(context.Background(), cacheKey, statusBytes, expires)
}

func getY2BUploadStatus(videoID string) (*Y2BVideoUploadStatus, bool, error) {
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
	var status Y2BVideoUploadStatus
	if err = json.Unmarshal(statusBytes, &status); err != nil {
		return nil, true, err // Exist, but cannot be unmarshalled
	}

	return &status, true, nil

}
