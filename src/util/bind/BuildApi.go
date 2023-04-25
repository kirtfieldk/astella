package bind

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/api"
	"github.com/kirtfieldk/astella/src/constants/routes"
	"github.com/kirtfieldk/astella/src/db"
	uuid "github.com/satori/go.uuid"
)

func injectDbConnectionIntoController(db *sql.DB) {

}

func BuildApi() {
	var dbConnection = db.CreateConnection()
	var baseHandler = api.NewBaseHandler(dbConnection)
	router := gin.Default()
	router.Use(requestIdMiddleware())

	router.GET(routes.GET_EVENT_BY_ID, baseHandler.GetEvent)
	router.POST(routes.GET_EVENT_BY_CITY, baseHandler.GetEventByCity)
	router.POST(routes.ADD_USER_TO_EVENT, baseHandler.AddUserToEvent)
	router.POST(routes.POST_MESSAGE_TO_EVENT, baseHandler.PostMessageToEvent)
	router.POST(routes.GET_MESSAGE_IN_EVENT, baseHandler.GetMessagesInEvent)
	router.POST(routes.LIKE_MESSAGE_IN_EVENT, baseHandler.UpvoteMessage)
	router.POST(routes.GET_USRS_LIKE_MESSAGE, baseHandler.GetUserUpvotes)

	router.Run("localhost:9000")
}
func requestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}
