package s3service

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadObject(bucket string, file *os.File, fileName string, session *s3.Client) {
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	_, err := session.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(fileName),
		Body:          fileBytes,
		ContentLength: *aws.Int64(size),
		ContentType:   aws.String(fileType),
	})

	if err != nil {
		log.Printf("\nFailed to upload object %v", err)
		return
	}

}

func GeneratePresignedUrl(bucket string, fileName string, session *s3.Client) (string, error) {
	var presigned string
	// svc := s3.New(session)
	// r, _ := svc.GetObjectRequest(&s3.GetObjectInput{
	// 	Bucket: aws.String("astellaapplicationmessages"),
	// 	Key:    aws.String("Screenshot 2022-07-06 at 8.53.17 PM.png"),
	// })
	// // Create the pre-signed url with an expiry
	// presigned, err := r.Presign(15 * time.Minute)
	// if err != nil {
	// 	fmt.Println("Failed to generate a pre-signed url: ", err)
	// 	return presigned, err
	// }

	// // Display the pre-signed url
	// fmt.Println("Pre-signed URL", presigned)
	return presigned, nil

}
