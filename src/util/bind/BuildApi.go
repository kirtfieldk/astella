package bind

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/api"
	"github.com/kirtfieldk/astella/src/constants/routes"
	"github.com/kirtfieldk/astella/src/db"
)

func injectDbConnectionIntoController(db *sql.DB) {

}

func BuildApi() {
	var dbConnection = db.CreateConnection()
	var baseHandler = api.NewBaseHandler(dbConnection)
	router := gin.Default()

	router.GET(routes.GET_EVENT_BY_ID, baseHandler.GetEvent)
	router.POST(routes.GET_EVENT_BY_CITY, baseHandler.GetEventByCity)
	router.POST(routes.ADD_USER_TO_EVENT, baseHandler.AddUserToEvent)
	router.POST(routes.POST_MESSAGE_TO_EVENT, baseHandler.PostMessageToEvent)

	router.Run("localhost:9000")
}
