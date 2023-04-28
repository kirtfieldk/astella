package userservice

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants/queries"
	eventservices "github.com/kirtfieldk/astella/src/services/eventServices"
	"github.com/kirtfieldk/astella/src/structures"
)

func MapUserRows(rows *sql.Rows) []structures.User {
	var users []structures.User
	for rows.Next() {
		var usr structures.User
		if err := rows.Scan(&usr.UUID, &usr.Username, &usr.Description, &usr.Created, &usr.Ig, &usr.Twitter, &usr.TikTok, &usr.AvatarUrl,
			&usr.ImgOne, &usr.ImgTwo, &usr.ImgThree); err != nil {
			log.Println("Issue mapping DB row for user")
		}

		users = append(users, usr)
	}
	return users
}

func UpdateUserProfile(user structures.User, conn *sql.DB) (bool, error) {
	_, err := conn.Exec(`UPDATE users SET username = $1, ig = $2, twitter = $3, toktok = $4, avatar_url = $5, 
		img_one = $5, img_two = $7, img_three = $8, description = $9 WHERE id = $10`, &user.Username,
		&user.Ig, &user.Twitter, &user.TikTok, &user.AvatarUrl, &user.ImgOne, &user.ImgTwo, &user.ImgThree, &user.Description)
	if err != nil {
		log.Panicln(err)
		return false, fmt.Errorf("Cannot Update User: " + user.UUID)
	}
	return true, nil
}

func GetUserEvents(userId string, conn *sql.DB) ([]structures.Event, error) {
	var events []structures.Event
	uuidInUrl, err := uuid.ParseBytes([]byte(userId))
	if err != nil {
		return nil, err
	}
	stmt, err := conn.Prepare(queries.GET_EVENTS_LOCATION_INFO_USER_IN)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return events, fmt.Errorf(`Unable to prepare statement to get user events.`)
	}
	rows, err := stmt.Query(uuidInUrl)
	if err != nil {
		log.Println(err)
		return events, fmt.Errorf(`Unable to query member's event.`)
	}
	events, err = eventservices.MapMultiLineRows(rows)
	if err != nil {
		return events, fmt.Errorf(`unable to map events`)
	}
	return events, nil
}
