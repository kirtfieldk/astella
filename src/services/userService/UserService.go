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
	"github.com/kirtfieldk/astella/src/util"
	uuidtransform "github.com/kirtfieldk/astella/src/util/uuidTransform"
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

func GetEventMembers(eventId string, page int, conn *sql.DB) (responses.UserListResponse, error) {
	var resp responses.UserListResponse
	eId, err := uuidtransform.StringToUuidTransform(eventId)
	if err != nil {
		return resp, fmt.Errorf("Failed to be UUID for User: " + eventId)
	}
	stmt, err := conn.Prepare(queries.GET_EVENT_USERS)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return resp, fmt.Errorf(`Unable to prepare statement to get user events.`)
	}
	rows, err := stmt.Query(eId, page*constants.LIMIT, constants.LIMIT)
	if err != nil {
		log.Println(err)
		return resp, fmt.Errorf(`Unable to query member's event.`)
	}
	users := MapUserRows(rows)
	resp.Data = users
	total, err := getTotalNumberOfUsersInEvent(eId, conn)
	if err != nil {
		return resp, err
	}
	resp.Info = structures.Info{
		Total: total,
		Count: len(users),
		Next:  (page*constants.LIMIT < total && len(users) != total),
	}
	return resp, nil
}

func GetEventUserIsMember(userId string, page int, conn *sql.DB) (responses.EventListResponse, error) {
	var resp responses.EventListResponse
	var total int
	uId, err := uuidtransform.StringToUuidTransform(userId)
	if err != nil {
		return resp, fmt.Errorf("Failed to be UUID for User: " + userId)
	}
	stmt, err := conn.Prepare(queries.GET_EVENTS_LOCATION_INFO_USER_IN)
	count, err := conn.Prepare(queries.GET_EVENTS_MEMBER_OF_COUNT)
	defer stmt.Close()
	defer count.Close()
	if err != nil {
		log.Println(err)
		return resp, fmt.Errorf(`Unable to prepare statement to get user events.`)
	}
	rows, err := stmt.Query(uId, util.CalcQueryStart(page), constants.LIMIT)
	if err != nil {
		log.Println(err)
		return resp, fmt.Errorf(`Unable to query member's event.`)
	}
	events, err := eventservices.MapMultiLineRows(rows)
	if err != nil {
		return resp, fmt.Errorf(`unable to map events`)
	}
	err = count.QueryRow(userId).Scan(&total)
	if err != nil {
		return resp, fmt.Errorf(`unable to map total`)
	}

	resp.Data = events
	resp.Info = structures.Info{
		Count: len(events),
		Total: total,
		Next:  page*constants.LIMIT < total && len(events) != total,
	}
	return resp, nil
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

func getTotalNumberOfUsersInEvent(eventId uuid.UUID, conn *sql.DB) (int, error) {
	var count int
	stmt, err := conn.Prepare(queries.GET_EVENTS_MEMBERS_COUNT)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return count, fmt.Errorf("Issue Here")
	}
	err = stmt.QueryRow(eventId).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
