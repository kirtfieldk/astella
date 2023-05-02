package util

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
)

func GetPageFromUrlQuery(c *gin.Context) int {
	page := c.DefaultQuery(constants.PAGE, "0")
	pagination, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Cannot convert to number " + page)
		pagination = 0
	}
	return pagination

}

func CalcQueryStart(page int) int {
	return page * constants.LIMIT
}
