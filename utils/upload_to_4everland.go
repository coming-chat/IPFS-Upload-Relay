package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
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

func Upload2ForeverLand(fileBytes []byte) (string, error) {

	// RandKey: file content hash

	fileHashBytes := sha256.Sum256(fileBytes)
	fileHash := hex.EncodeToString(fileHashBytes[:])

	log.Println("New file upload request with hash: ", fileHash)

	svc, err := prepare()
	if err != nil {
		return "", err
	}

	// Prepare CID
	var cid string

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
				Body:   bytes.NewReader(fileBytes),
				Bucket: &global.ForeverLand_Bucket,
				Key:    aws.String(fileHash),
			})
			if err != nil {
				return "", err
			} else {
				// cid = *uploadResp.ETag // Unable to handle here, need another head
				headResp, _ := svc.HeadObject(&s3.HeadObjectInput{
					Bucket: &global.ForeverLand_Bucket,
					Key:    aws.String(fileHash),
				})
				cid = *headResp.ETag
			}
		default:
			return "", err
		}
	} else {
		cid = *headResp.ETag
	}

	return strings.ReplaceAll(cid, "\"", ""), nil

}
