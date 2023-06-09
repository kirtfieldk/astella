package eventservices

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	b64 "encoding/base64"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/api/responses"
	"github.com/kirtfieldk/astella/src/constants"
	"github.com/kirtfieldk/astella/src/constants/queries"
	locationservice "github.com/kirtfieldk/astella/src/services/locationService"
	"github.com/kirtfieldk/astella/src/structures"
	"github.com/kirtfieldk/astella/src/util"
	uuidtransform "github.com/kirtfieldk/astella/src/util/uuidTransform"
)

// We can easily create events > FIx.
func CreateEvent(eventInfo structures.Event, conn *sql.DB) (responses.EventListResponse, error) {
	var resp responses.EventListResponse
	tx, err := conn.Begin()
	if err != nil {
		return resp, err
	}
	locationInfoDb, err := locationservice.CreateLocationInfo(eventInfo.LocationInfo, tx)
	if err != nil {
		return resp, err
	}
	eventId, err := insertEventIntoDb(eventInfo, locationInfoDb.Id, tx)
	if err != nil {
		return resp, err
	}
	userId, eventIdUuid, err := uuidtransform.ParseTwoIds(eventInfo.UserId, eventId)
	if err != nil {
		log.Printf("Issue parsing ids: %v AND %v\n", userId, eventIdUuid)
	}
	_, err = addUserToEventTransaction(userId, eventIdUuid, tx)
	if err != nil {
		return resp, err
	}
	err = tx.Commit()
	if err != nil {
		return resp, fmt.Errorf("Error with event transaction %v", err)
	}
	event, err := GetEvent(eventId, conn)
	if err != nil {
		return resp, err
	}
	resp.Data = append(resp.Data, event)
	resp.Info.Count = 1
	resp.Info.Total = 1
	resp.Info.Page = 0
	resp.Info.Next = false

	return resp, nil
}

func GetEvent(id string, conn *sql.DB) (structures.Event, error) {
	var event structures.Event
	uuidInUrl, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return event, err
	}

	stmt, err := conn.Prepare(queries.GET_EVENT_BY_ID_AND_LOCATION_INFO)
	defer stmt.Close()
	if err != nil {
		return event, err
	}
	row := stmt.QueryRow(uuidInUrl, time.Now().UTC())

	event, err = mapSingleRowQuery(row)
	if err != nil {
		log.Println(err)
		return event, fmt.Errorf("No Event with UUID:" + id)
	}
	return event, nil
}

func GetEventsByCity(city string, page int, conn *sql.DB) (responses.EventListResponse, error) {
	var response responses.EventListResponse
	if err := conn.Ping(); err != nil {
		return response, err
	}
	stmt, err := conn.Prepare(queries.GET_EVENT_BY_CITY_AND_LOCATION_INFO)
	defer stmt.Close()
	if err != nil {
		return response, err
	}
	rows, err := stmt.Query(city, time.Now().UTC(), util.CalcQueryStart(page), constants.LIMIT)
	defer rows.Close()
	events, err := MapMultiLineRows(rows)
	if err != nil {
		return response, getErrorMessage(err, city)
	}
	total, err := getTotalCountOfEventsInCity(city, conn)
	if err != nil {
		return response, fmt.Errorf("Issue Here")
	}
	response.Data = events
	response.Info = structures.Info{
		Count: len(events),
		Total: total,
		Page:  page,
		Next:  page*constants.LIMIT < total && len(events) != total,
	}

	return response, nil
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
	event, err := mapSingleRowQuery(conn.QueryRow(queries.GET_EVENT_BY_ID_AND_LOCATION_INFO, eventId, time.Now().UTC()))
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
		if event.IsPublic || (!event.IsPublic && string(decodedCode) == code) {
			return addUserToEvent(uId, eId, conn)
		}
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
	if err := row.Scan(&event.Id, &event.Name, &event.Created, &event.Description, &event.IsPublic, &event.Code,
		&location.Id, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
		&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.City,
		&event.UserId, &event.EndTime); err != nil {
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
			&event.UserId, &event.EndTime); err != nil {
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
	_, err := conn.Exec(queries.ADD_USER_TO_EVENT, userId, eventId, time.Now().UTC())
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func addUserToEventTransaction(userId uuid.UUID, eventId uuid.UUID, conn *sql.Tx) (bool, error) {
	_, err := conn.Exec(queries.ADD_USER_TO_EVENT, userId, eventId, time.Now().UTC())
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

// We can easily create events > FIx.
func insertEventIntoDb(eventInfo structures.Event, lId string, conn *sql.Tx) (string, error) {
	var id string
	encodedCode := b64.StdEncoding.EncodeToString([]byte(eventInfo.Code))
	stmt, err := conn.Prepare(queries.INSERT_EVENT_INTO_DB)
	defer stmt.Close()
	if err != nil {
		return id, err
	}
	err = stmt.QueryRow(&eventInfo.Name, time.Now().UTC(), &eventInfo.Description,
		&eventInfo.IsPublic, &encodedCode, &lId, &eventInfo.Duration,
		time.Now().UTC().Add(time.Hour*time.Duration(eventInfo.Duration)), &eventInfo.UserId).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("Unable to create event %v", err)
	}
	if err != nil {
		return id, fmt.Errorf("Unable to fetch event")
	}
	return id, nil
}

func getTotalCountOfEventsInCity(city string, conn *sql.DB) (int, error) {
	stmt, err := conn.Prepare(queries.GET_EVENT_BY_CITY_AND_LOCATION_INFO_COUNT)
	var count int
	defer stmt.Close()
	if err != nil {
		log.Panicln(err)
		log.Println("Issue formating query for total number of events in city: " + city)
		return count, err
	}
	err = stmt.QueryRow(city, time.Now().UTC()).Scan(&count)
	if err != nil {
		log.Println("Issue getting total number of events in city: " + city)
		return count, err
	}
	return count, nil
}
