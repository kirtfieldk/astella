package messageservice

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants/queries"
	locationservice "github.com/kirtfieldk/astella/src/services/locationService"
	"github.com/kirtfieldk/astella/src/structures"
	uuidtransform "github.com/kirtfieldk/astella/src/util/uuidTransform"
)

func getMessages(eventId string, conn *sql.DB) ([]structures.Message, error) {
	uuidInUrl, err := uuidtransform.StringToUuidTransform(eventId)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(queries.GET_MESSAGES_IN_EVENT, uuidInUrl)
	defer rows.Close()
	if err != nil {
		return nil, getErrorMessageForNoMessages(err, eventId)
	}
	return mapRowsToMessages(rows), nil

}

func PostMessage(msg structures.Message, conn *sql.DB) (bool, error) {
	eventId, err := uuidtransform.StringToUuidTransform(msg.EventId)
	if err != nil {
		log.Printf("? Failed to be UUID for Event: " + msg.EventId)
		return false, err
	}
	userId, err := uuidtransform.StringToUuidTransform(msg.UserId)
	if err != nil {
		log.Printf("Failed to be UUID for User: " + msg.UserId)
		return false, err
	}
	if isRequestInArea(eventId, msg.Latitude, msg.Longitude, conn) && isUserInEvent(userId, eventId, conn) {
		return insertMessageIntoDB(msg, conn), nil
	}
	return false, nil

}

func getErrorMessageForNoMessages(err error, id string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("Messages For EventId %s: no messages", id)
	}
	return fmt.Errorf("MessagesEventById %s: %v", id, err)
}

func mapRowsToMessages(rows *sql.Rows) []structures.Message {
	var messages []structures.Message
	for rows.Next() {
		var msg structures.Message
		if err := rows.Scan(&msg.UUID, &msg.Content, &msg.Created, &msg.EventId, &msg.ParentId, &msg.Upvotes,
			&msg.Pinned, &msg.Latitude, &msg.Longitude); err != nil {
			//continue
		}

		messages = append(messages, msg)
	}
	return messages
}

func isRequestInArea(id uuid.UUID, lat float32, long float32, conn *sql.DB) bool {
	log.Println(id)
	location := locationservice.GetEventLocation(id, conn)
	topLeftPt := locationservice.CreatePoint(location.TopLeftLat, location.TopLeftLon)
	topRightPt := locationservice.CreatePoint(location.TopRightLat, location.TopRightLon)
	bottomLeftPt := locationservice.CreatePoint(location.BottomLeftLat, location.BottomLeftLon)
	bottomRightPt := locationservice.CreatePoint(location.BottomLeftLat, location.BottomRightLon)
	return locationservice.CheckPointInArea(lat, long, topLeftPt, topRightPt, bottomLeftPt, bottomRightPt)
}

func isUserInEvent(userId uuid.UUID, eventId uuid.UUID, conn *sql.DB) bool {
	stmt, err := conn.Prepare(queries.FIND_USER_IN_EVENT)
	if err != nil {
		log.Println("Touble preparing statement for isUserInEvent")
	}
	_, err = stmt.Query(userId, eventId)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func insertMessageIntoDB(msg structures.Message, conn *sql.DB) bool {
	if msg.ParentId != "" {
		_, err := conn.Exec(queries.INSERT_MESSAGE_WITH_PARENT_ID, &msg.Content, &msg.UserId, time.Now(), &msg.EventId, &msg.ParentId, &msg.Upvotes, &msg.Pinned,
			&msg.Latitude, &msg.Longitude)
		if err != nil {
			log.Println(err)
			return false
		}
	} else {
		_, err := conn.Exec(queries.INSERT_MESSAGE_WITHOUT_PARENT_ID, &msg.Content, &msg.UserId, time.Now(), &msg.EventId, &msg.Upvotes, &msg.Pinned,
			&msg.Latitude, &msg.Longitude)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}
