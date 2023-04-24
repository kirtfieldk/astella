package services

import (
	"database/sql"
	"fmt"

	_ "github.com/kirtfieldk/astella/src/db"
	_ "github.com/kirtfieldk/astella/src/structures"
)

// func getMessages(eventId string) ([]structures.Message, error) {
// 	var connection *sql.DB
// 	connection = db.GetConnection()

// 	rows, err := connection.Query("SELECT * FROM messages WHERE event_id = ? ORDER BY created ASC LIMIT 30", eventId)
// 	if err != nil {
// 		getErrorMessageForNoMessages(err, eventId)
// 	}

// }

func getErrorMessageForNoMessages(err error, id string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("Messages For EventId %s: no messages", id)
	}
	return fmt.Errorf("MessagesEventById %s: %v", id, err)
}
