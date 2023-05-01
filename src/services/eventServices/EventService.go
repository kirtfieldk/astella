package eventservices

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	b64 "encoding/base64"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants/queries"
	locationservice "github.com/kirtfieldk/astella/src/services/locationService"
	"github.com/kirtfieldk/astella/src/structures"
	uuidtransform "github.com/kirtfieldk/astella/src/util/uuidTransform"
)

// We can easily create events > FIx.
func CreateEvent(eventInfo structures.Event, conn *sql.DB) (bool, error) {
	tx, err := conn.Begin()
	if err != nil {
		return false, err
	}

	locationId, err := locationservice.CreateLocationInfo(eventInfo.LocationInfo, tx)
	if err != nil {
		return false, err
	}
	if _, err := insertEventIntoDb(eventInfo, locationId, tx); err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("Error with event transaction %v", err)
	}
	return true, nil
}

func GetEvent(id string, conn *sql.DB) (*structures.Event, error) {
	uuidInUrl, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return nil, err
	}

	stmt, err := conn.Prepare(queries.GET_EVENT_BY_ID_AND_LOCATION_INFO)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(uuidInUrl)

	event, err := mapSingleRowQuery(row)
	if err != nil {
		log.Println(err)
		return &event, fmt.Errorf("No Event with UUID:" + id)
	}
	date, err := time.Parse(time.RFC3339, event.EndTime)
	if err != nil {
		log.Println(err)
		log.Println(event.EndTime)
		return &event, fmt.Errorf(`Issue Parsing date %s`, event.EndTime)
	}
	if event.Expired || time.Now().UTC().After(date) {
		expireEvent(event, conn)
		return &event, fmt.Errorf("Event has expired.")
	}
	return &event, nil
}

func GetEventsByCity(city string, conn *sql.DB) ([]structures.Event, error) {
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	stmt, err := conn.Prepare(queries.GET_EVENT_BY_CITY_AND_LOCATION_INFO)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	log.Println(city)
	rows, err := stmt.Query(city)
	defer rows.Close()
	events, err := MapMultiLineRows(rows)
	if err != nil {
		return events, getErrorMessage(err, city)
	}
	for _, e := range events {
		date, err := time.Parse(time.RFC3339, e.EndTime)
		if err != nil {
			log.Println(`Issue Parsing date ` + e.EndTime)
		}
		if time.Now().UTC().After(date) {
			expireEvent(e, conn)
			log.Println(`Issue Parsing date ` + e.Id)
		}
	}
	return events, nil
}

func AddUserToEvent(code string, userId string, eventId string, cords structures.Point, conn *sql.DB) (bool, error) {
	eId, err := uuidtransform.StringToUuidTransform(eventId)
	if err != nil {
		return false, fmt.Errorf("Failed to be UUID for Event: " + eventId)
	}
	uId, err := uuidtransform.StringToUuidTransform(userId)
	if err != nil {
		return false, fmt.Errorf("Failed to be UUID for User: " + userId)
	}
	event, err := mapSingleRowQuery(conn.QueryRow(queries.GET_EVENT_BY_ID_AND_LOCATION_INFO, eventId))
	if err != nil {
		log.Println(err)
		return false, getErrorMessage(err, eventId)
	}
	decodedCode, err := b64.StdEncoding.DecodeString(event.Code)
	if err != nil {
		log.Panicln("unable to decode code: " + event.Code)
	}
	if checkPointInEvent(cords.Latitude, cords.Longitude,
		structures.Point{Latitude: event.LocationInfo.TopLeftLat, Longitude: event.LocationInfo.TopLeftLon},
		structures.Point{Latitude: event.LocationInfo.TopRightLat, Longitude: event.LocationInfo.TopRightLon},
		structures.Point{Latitude: event.LocationInfo.BottomLeftLat, Longitude: event.LocationInfo.BottomLeftLon},
		structures.Point{Latitude: event.LocationInfo.BottomRightLat, Longitude: event.LocationInfo.BottomRightLon}) {
		if event.Expired {
			return false, fmt.Errorf("Event has expired.")
		}
		if event.IsPublic || (!event.IsPublic && string(decodedCode) == code) {
			return addUserToEvent(uId, eId, conn)
		}
	}
	return false, fmt.Errorf("Cannot add user to event")

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
	if err := row.Scan(&event.Id, &event.Name, &event.Created, &event.Description, &event.IsPublic, &event.Code,
		&location.Id, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
		&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.City,
		&event.Expired, &event.EndTime); err != nil {
		return event, err
	}
	event.LocationInfo = location
	return event, nil
}

func MapMultiLineRows(rows *sql.Rows) ([]structures.Event, error) {
	var events []structures.Event = make([]structures.Event, 0)
	for rows.Next() {
		var event structures.Event
		var location structures.LocationInfo
		if err := rows.Scan(&event.Id, &event.Name, &event.Created, &event.Description, &event.IsPublic, &event.Code,
			&location.Id, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
			&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.City,
			&event.Expired, &event.EndTime); err != nil {
			//continue
			log.Println(err)
		}
		event.LocationInfo = location
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

func addUserToEvent(userId uuid.UUID, eventId uuid.UUID, conn *sql.DB) (bool, error) {
	_, err := conn.Exec("INSERT INTO members (user_id, event_id, created) values ($1, $2, $3)", userId, eventId, time.Now())
	if err != nil {
		return false, err
	}
	return true, nil
}

// We can easily create events > FIx.
func insertEventIntoDb(eventInfo structures.Event, lId uuid.UUID, conn *sql.Tx) (bool, error) {
	encodedCode := b64.StdEncoding.EncodeToString([]byte(eventInfo.Code))
	stmt, err := conn.Prepare(queries.INSERT_EVENT_INTO_DB)
	defer stmt.Close()
	if err != nil {
		return false, err
	}
	log.Println(lId)
	_, err = stmt.Exec(&eventInfo.Name, time.Now().UTC(), &eventInfo.Description,
		&eventInfo.IsPublic, &encodedCode, &lId, &eventInfo.Duration,
		time.Now().UTC().Add(time.Hour*time.Duration(eventInfo.Duration)), false)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf("Unable to create event %v", err)
	}
	return true, nil
}

func expireEvent(event structures.Event, conn *sql.DB) {
	_, err := conn.Exec(queries.EXPIRE_EVENT, event.Id)
	if err != nil {
		log.Println(`Issue expiring event with UUID: `, event.Id)
	}
}
