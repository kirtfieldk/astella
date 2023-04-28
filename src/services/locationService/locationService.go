package locationservice

import (
	"database/sql"
	"fmt"
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
		log.Println(id)
	}
	return locationInfo

}

func CheckPointInArea(lat float32, long float32, topLeft structures.Point, topRight structures.Point, bottomLeft structures.Point, bottomRight structures.Point) bool {
	var minX = math.Min(math.Min(float64(bottomLeft.Latitude), float64(topLeft.Latitude)), math.Min(float64(bottomRight.Latitude), float64(topRight.Latitude)))
	var maxX = math.Max(math.Max(float64(bottomLeft.Latitude), float64(topLeft.Latitude)), math.Max(float64(bottomRight.Latitude), float64(topRight.Latitude)))
	var minY = math.Min(math.Min(float64(bottomLeft.Longitude), float64(topLeft.Longitude)), math.Min(float64(bottomRight.Longitude), float64(topRight.Longitude)))
	var maxY = math.Max(math.Max(float64(bottomLeft.Longitude), float64(topLeft.Longitude)), math.Max(float64(bottomRight.Longitude), float64(topRight.Longitude)))
	var updatedLat = float64(lat)
	var updatedLon = float64(long)

	return minX <= updatedLat && updatedLat <= maxX && minY <= updatedLon && updatedLon <= maxY

}

func CreatePoint(lat float32, long float32) structures.Point {
	return structures.Point{Latitude: lat, Longitude: long}
}

func CreateLocationInfo(locationInfo structures.LocationInfo, tx *sql.Tx) (uuid.UUID, error) {
	var id uuid.UUID
	stmt, err := tx.Prepare(queries.INSERT_LOCATION_INTO_DB_RETURN_ID)
	if err != nil {
		log.Println(err)
		return id, fmt.Errorf("Error preparing insert location stmt")
	}
	defer stmt.Close()
	err = stmt.QueryRow(&locationInfo.TopLeftLat, &locationInfo.TopLeftLon, &locationInfo.TopRightLat, &locationInfo.TopRightLon,
		&locationInfo.BottomLeftLat, &locationInfo.BottomLeftLon, &locationInfo.BottomRightLat, &locationInfo.BottomRightLon, &locationInfo.City).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("Error mapping Id from LocationInfo create stmt")
	}
	return id, err
}
