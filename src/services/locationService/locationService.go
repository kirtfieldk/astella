package locationservice

import (
	"database/sql"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants/queries"
	"github.com/kirtfieldk/astella/src/structures"
)

func GetEventLocation(id uuid.UUID, conn *sql.DB) structures.LocationInfo {
	var locationInfo structures.LocationInfo
	s := conn.QueryRow(queries.GET_LOCATION_FOR_EVENT, id)
	err := s.Scan(&locationInfo.UUID, &locationInfo.TopLeftLat, &locationInfo.TopLeftLon,
		&locationInfo.TopRightLat, &locationInfo.TopRightLon, &locationInfo.BottomLeftLat,
		&locationInfo.BottomLeftLon, &locationInfo.BottomRightLat, &locationInfo.BottomRightLon)
	if err != nil {
		log.Println(err)
	}
	return locationInfo

}

func CheckPointInArea(lat float32, long float32, topLeft structures.Point, topRight structures.Point, bottomLeft structures.Point, bottomRight structures.Point) bool {
	var minX = math.Min(float64(topLeft.Latitude), float64(topRight.Latitude))
	var maxX = math.Max(float64(topLeft.Latitude), float64(topRight.Latitude))
	var minY = math.Min(float64(bottomLeft.Longitude), float64(topLeft.Longitude))
	var maxY = math.Max(float64(bottomLeft.Longitude), float64(topLeft.Longitude))
	var updatedLat = float64(lat)
	var updatedLon = float64(long)

	return minX <= updatedLat && updatedLat <= maxX && minY <= updatedLon && updatedLon <= maxY

}

func CreatePoint(lat float32, long float32) structures.Point {
	return structures.Point{Latitude: lat, Longitude: long}
}
