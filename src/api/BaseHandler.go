package api

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws/session"
)

type BaseHandler struct {
	DB         *sql.DB
	AwsSession *session.Session
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB, awsSession *session.Session) *BaseHandler {
	return &BaseHandler{
		DB:         db,
		AwsSession: awsSession,
	}
}
