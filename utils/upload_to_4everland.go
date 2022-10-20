package utils

import (
	"fmt"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	sess *session.Session
)

func prepare() (*s3.S3, error) {
	if sess == nil {
		var err error
		sess, err = session.NewSession(&aws.Config{
			Endpoint: aws.String("https://endpoint.4everland.co/"),
			Region:   aws.String("us-west-1"),
		})
		if err != nil {
			return nil, err
		}
	}

	return s3.New(sess), nil
}

func Upload2ForeverLand(r io.ReadSeeker) (string, int64, error) {

	// RandKey: file content hash

	fileHash := CalcFileHash(r)

	log.Println("New file upload request with hash: ", fileHash)

	svc, err := prepare()
	if err != nil {
		return "", -1, err
	}

	// Prepare CID
	var (
		cid      string
		filesize int64
	)

	// Check if already exists
	headResp, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &global.ForeverLand_Bucket,
		Key:    aws.String(fileHash),
	})
	if err != nil {
		switch err.(awserr.Error).Code() {
		case "NotFound":
			// Upload file
			_, err := svc.PutObject(&s3.PutObjectInput{
				Body:   r,
				Bucket: &global.ForeverLand_Bucket,
				Key:    aws.String(fileHash),
			})
			if err != nil {
				return "", -1, err
			} else {
				// cid = *uploadResp.ETag // Unable to handle here, need another head
				headResp, _ := svc.HeadObject(&s3.HeadObjectInput{
					Bucket: &global.ForeverLand_Bucket,
					Key:    aws.String(fileHash),
				})
				cid = *headResp.ETag
				filesize = *headResp.ContentLength
			}
		default:
			return "", -1, err
		}
	} else {
		cid = *headResp.ETag
		filesize = *headResp.ContentLength
	}

	// Request once to ensure file pinned
	go (&http.Client{}).Get(fmt.Sprintf("https://4everland.io/ipfs/%s", cid))

	return strings.ReplaceAll(cid, "\"", ""), filesize, nil

}
