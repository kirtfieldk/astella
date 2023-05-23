package s3service

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadObject(bucket string, file *os.File, fileName string, session *session.Session) (string, error) {
	uploader := s3manager.NewUploader(session)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		log.Printf("\nFailed to upload object %v", err)
		return fileName, err
	}
	log.Printf("\nSuccessfully upload %q to %q", fileName, bucket)
	return output.Location, nil
}
