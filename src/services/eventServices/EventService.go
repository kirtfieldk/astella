package eventservices

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants/queries"
	locationservice "github.com/kirtfieldk/astella/src/services/locationService"
	"github.com/kirtfieldk/astella/src/structures"
)

// We can easily create events > FIx.
func CreateEvent(eventInfo structures.Event, conn *sql.DB) (bool, error) {
	_, err := conn.Exec(`Insert INTO events (event_name, created, description, top_left, top_right, bottom_left, bottom_right,
    	public, code) VALUES (?,?,?, point(?, ?), point(?, ?), point(?, ?),point(?, ?) ,?, ?)`,
		&eventInfo.Name, &eventInfo.Created, &eventInfo.Description, &eventInfo.Public, &eventInfo.Code)
	if err != nil {
		return false, fmt.Errorf("Unable to create event %v", err)
	}

	return true, nil
}

func GetEvent(id string, conn *sql.DB) (*structures.Event, error) {
	uuidInUrl, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return nil, err
	}

	stmt, err := conn.Prepare(queries.GET_EVENT_BY_ID_AND_LOCATION_INFO)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(uuidInUrl)

	event, err := mapSingleRowQuery(row)
	if err != nil {
		log.Println(err)
		return &event, err
	}
	return &event, nil
}

func GetEventsByCity(city string, conn *sql.DB) ([]structures.Event, error) {
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	stmt, err := conn.Prepare(queries.GET_EVENT_BY_CITY_AND_LOCATION_INFO)
	if err != nil {
		return nil, err
	}
	log.Println(city)
	rows, err := stmt.Query(city)
	defer rows.Close()
	events, err := mapMultiLineRows(rows)
	if err != nil {
		log.Println(err)
		return events, getErrorMessage(err, city)
	}
	log.Println(events)
	return events, nil
}

func AddUserToEvent(code string, userId string, eventId string, conn *sql.DB) (bool, error) {
	row := conn.QueryRow("Select * from Event where UUID = $1", eventId)
	event, err := mapSingleRowQuery(row)
	if err != nil {
		return false, getErrorMessage(err, eventId)
	}
	if event.Public {
		// adduser
	} else if event.Code == code {

	}
	return false, nil

}

func DeleteEvent(id string, conn *sql.DB) (bool, error) {
	_, err := conn.Exec("Delete from Event where UUID = $1", id)
	if err != nil {
		return false, getErrorMessage(err, id)
	}
	return true, nil
}

func IsUserInEvent(userId uuid.UUID, eventId uuid.UUID) {

}

func mapSingleRowQuery(row *sql.Row) (structures.Event, error) {
	var event structures.Event
	var location structures.LocationInfo
	if err := row.Scan(&event.UUID, &event.Name, &event.Created, &event.Description, &event.Public, &event.Code,
		&location.UUID, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
		&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.City); err != nil {
		return event, err
	}
	event.Location = location
	return event, nil
}

func mapMultiLineRows(rows *sql.Rows) ([]structures.Event, error) {
	var events []structures.Event
	for rows.Next() {
		var event structures.Event
		var location structures.LocationInfo
		if err := rows.Scan(&event.UUID, &event.Name, &event.Created, &event.Description, &event.Public, &event.Code,
			&location.UUID, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
			&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.City); err != nil {
			//continue
			log.Println(err)
		}

		event.Location = location
		events = append(events, event)

	}
	return events, nil
}

// topLeft.x <= pt < topRight (take min), same for Y
func checkPointInEvent(lat float32, long float32, topLeft structures.Point, topRight structures.Point, bottomLeft structures.Point, bottomRight structures.Point) bool {
	return locationservice.CheckPointInArea(lat, long, topLeft, topRight, bottomLeft, bottomRight)
}

func getErrorMessage(err error, id string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("EventById %s: no such Event %v", id, err)
	}
	return fmt.Errorf("EventById %s: %v", id, err)
}