package storage

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadFile func
func UploadFile(uploadFileDir string, hashedName string) error {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION")),
	})

	if err != nil {
		return err
	}

	upFile, err := os.Open(uploadFileDir)

	if err != nil {
		return err
	}

	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()

	var fileSize int64 = upFileInfo.Size()

	fileBuffer := make([]byte, fileSize)

	upFile.Read(fileBuffer)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:                  aws.String(hashedName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	return err
}
