package api

import "database/sql"

type BaseHandler struct {
	DB *sql.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		DB: db,
	}
}
