package utils

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"strings"
)

var (
	awsS3Session *session.Session
)

func awsS3prepare() (*s3.S3, error) {
	if awsS3Session == nil {
		var err error
		awsS3Session, err = session.NewSession(&aws.Config{
			Endpoint: aws.String("https://s3.us-east-2.amazonaws.com/coming-upload/"),
			Region:   aws.String("us-east-2"),
			LogLevel: aws.LogLevel(aws.LogDebugWithSigning),
		})
		if err != nil {
			return nil, err
		}
	}

	return s3.New(awsS3Session), nil
}

func UploadToAwsS3(r io.ReadSeeker, contentType string) (string, int64, error) {

	// RandKey: file content hash
	fileData, err := io.ReadAll(r)
	if err != nil {
		return "", 0, err
	}
	_, _ = r.Seek(0, io.SeekStart)
	cid, err := GetIPFSCid(fileData)
	if err != nil {
		return "", 0, err
	}

	log.Println("New file upload request with Cid: ", cid)

	svc, err := awsS3prepare()
	if err != nil {
		return "", 0, err
	}

	// Prepare CID
	var (
		filesize int64
	)

	// Check if already exists
	headResp, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &global.ForeverLand_Bucket,
		Key:    aws.String(cid),
	})
	if err != nil {
		switch err.(awserr.Error).Code() {
		case "NotFound":
			// Upload file
			_, err := svc.PutObject(&s3.PutObjectInput{
				Body:        r,
				Bucket:      &global.AwsS3_Bucket,
				Key:         aws.String(cid),
				ContentType: &contentType,
			})
			if err != nil {
				return "", 0, err
			} else {
				headResp, _ := svc.HeadObject(&s3.HeadObjectInput{
					Bucket: &global.ForeverLand_Bucket,
					Key:    aws.String(cid),
				})
				filesize = *headResp.ContentLength
			}
		default:
			return "", 0, err
		}
	} else {
		filesize = *headResp.ContentLength
	}

	return strings.ReplaceAll(cid, "\"", ""), filesize, nil

}
