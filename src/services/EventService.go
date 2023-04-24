package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants"
	"github.com/kirtfieldk/astella/src/db"
	"github.com/kirtfieldk/astella/src/structures"
)

var (
	ctx context.Context
	con *sql.DB
)

type EventService interface {
	GetEvent(id string) ([]structures.Event, error)
	getEvents(lat float32, lon float32) ([]structures.Event, error)
	createEvent(eventInfo structures.Event) (bool, error)
	deleteEvent(id string) (bool, error)
	addUserToEvent(code string, userId string, eventId string) (bool, error)
}

// We can easily create events > FIx.
func createEvent(eventInfo structures.Event) (bool, error) {
	_, err := con.ExecContext(ctx, `Insert INTO events (event_name, created, description, top_left, top_right, bottom_left, bottom_right,
    	public, code) VALUES (?,?,?, point(?, ?), point(?, ?), point(?, ?),point(?, ?) ,?, ?)`,
		&eventInfo.Name, &eventInfo.Created, &eventInfo.Description, &eventInfo.Public, &eventInfo.Code)
	if err != nil {
		return false, fmt.Errorf("Unable to create event %v", err)
	}

	return true, nil
}

func GetEvent(id string) (*structures.Event, error) {
	uuidInUrl, err := uuid.ParseBytes([]byte(id))
	if err != nil {
		return nil, err
	}
	if err := db.DbConnection.Ping(); err != nil {
		return nil, err
	}
	stmt, err := db.DbConnection.Prepare(constants.GET_EVENT_BY_ID_AND_LOCATION_INFO)
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

func GetEventsByCity(city string) ([]structures.Event, error) {
	log.Println(city)
	if err := db.DbConnection.Ping(); err != nil {
		return nil, err
	}
	stmt, err := db.DbConnection.Prepare(constants.GET_EVENT_BY_CITY_AND_LOCATION_INFO)
	if err != nil {
		return nil, err
	}
	row, err := stmt.Query(city)
	events, err := mapMultiLineRows(row)
	if err != nil {
		return events, getErrorMessage(err, city)
	}
	return events, nil
}

func AddUserToEvent(code string, userId string, eventId string) (bool, error) {
	// con := db.GetConnection()
	row := con.QueryRow("Select * from Event where UUID = $1", eventId)
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

func DeleteEvent(id string) (bool, error) {
	// con := db.GetConnection()
	_, err := con.Exec("Delete from Event where UUID = $1", id)
	if err != nil {
		return false, getErrorMessage(err, id)
	}
	return true, nil

}

func mapSingleRowQuery(row *sql.Row) (structures.Event, error) {
	var event structures.Event
	var location structures.LocationInfo
	if err := row.Scan(&event.UUID, &event.Name, &event.Created, &event.Description, &event.Public, &event.Code,
		&location.UUID, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
		&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.Latitude,
		&location.Longitude, &location.City); err != nil {
		return event, err
	}
	event.Location = location
	return event, nil
}

func mapMultiLineRows(rows *sql.Rows) ([]structures.Event, error) {
	var events []structures.Event
	defer rows.Close()
	for rows.Next() {
		var event structures.Event
		var location structures.LocationInfo
		if err := rows.Scan(&event.UUID, &event.Name, &event.Created, &event.Description, &event.Public, &event.Code,
			&location.UUID, &location.TopLeftLat, &location.TopLeftLon, &location.TopRightLat, &location.TopRightLon,
			&location.BottomRightLat, &location.BottomRightLon, &location.BottomLeftLat, &location.BottomLeftLon, &location.Latitude,
			&location.Longitude, &location.City); err != nil {
			//continue
		}

		event.Location = location
		events = append(events, event)

	}
	return events, nil
}

func getErrorMessage(err error, id string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("EventById %s: no such Event %v", id, err)
	}
	return fmt.Errorf("EventById %s: %v", id, err)
}
