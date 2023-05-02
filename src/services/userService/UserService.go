package userservice

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/api/responses"
	"github.com/kirtfieldk/astella/src/constants"
	"github.com/kirtfieldk/astella/src/constants/queries"
	eventservices "github.com/kirtfieldk/astella/src/services/eventServices"
	"github.com/kirtfieldk/astella/src/structures"
)

func MapUserRows(rows *sql.Rows) []structures.User {
	var users []structures.User
	for rows.Next() {
		var usr structures.User
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.Description, &usr.Created, &usr.Ig, &usr.Twitter, &usr.TikTok, &usr.AvatarUrl,
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
		return false, fmt.Errorf("Cannot Update User: " + user.Id)
	}
	return true, nil
}

func GetUserEvents(userId string, page int, conn *sql.DB) (responses.EventListResponse, error) {
	var response responses.EventListResponse
	uuidInUrl, err := uuid.ParseBytes([]byte(userId))
	if err != nil {
		return response, err
	}
	events, err := getEventUserIsMember(uuidInUrl, page, conn)
	if err != nil {
		log.Println(err)
		return response, fmt.Errorf("Unable to get users")
	}
	total, err := getTotalNumberOfEventsUserIsApartOf(uuidInUrl, conn)
	if err != nil {
		return response, err
	}
	response.Data = events
	response.Info = structures.Info{
		Count: len(events),
		Total: total,
		Page:  page,
	}
	return response, nil
}

func getEventUserIsMember(userId uuid.UUID, page int, conn *sql.DB) ([]structures.Event, error) {
	var events []structures.Event
	stmt, err := conn.Prepare(queries.GET_EVENTS_LOCATION_INFO_USER_IN)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return events, fmt.Errorf(`Unable to prepare statement to get user events.`)
	}
	rows, err := stmt.Query(userId, page*constants.LIMIT, constants.LIMIT)
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

func getTotalNumberOfEventsUserIsApartOf(userId uuid.UUID, conn *sql.DB) (int, error) {
	var count int
	stmt, err := conn.Prepare(queries.GET_EVENTS_LOCATION_INFO_USER_IN_COUNT)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return count, fmt.Errorf("Issue Here")
	}
	err = stmt.QueryRow(userId).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
