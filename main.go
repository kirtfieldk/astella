package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/api"
	"github.com/kirtfieldk/astella/src/constants"
	"github.com/kirtfieldk/astella/src/db"
)

func main() {
	router := gin.Default()
	router.GET(constants.GET_EVENT_BY_ID, api.GetEvent)
	router.POST(constants.GET_EVENT_BY_CITY, api.GetEventByCity)
	router.POST(constants.ADD_USER_TO_EVENT, api.AddUserToEvent)
	db.CreateConnection()
	// defer db.DbConnection.Close()

	router.Run("localhost:9000")
}
