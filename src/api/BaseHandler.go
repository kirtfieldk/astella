package api

import (
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BaseHandler struct {
	DB        *sql.DB
	S3Session *s3.Client
	AwsBucket string
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB, awsSession *s3.Client, bucket string) *BaseHandler {
	return &BaseHandler{
		DB:        db,
		S3Session: awsSession,
		AwsBucket: bucket,
	}
}
