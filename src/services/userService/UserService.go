package userservice

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

func mapUserRow(rows *sql.Row) structures.User {
	var usr structures.User
	if err := rows.Scan(&usr.Id, &usr.Username, &usr.Description, &usr.Created, &usr.Ig, &usr.Twitter, &usr.TikTok, &usr.AvatarUrl,
		&usr.ImgOne, &usr.ImgTwo, &usr.ImgThree); err != nil {
		log.Println("Issue mapping DB row for user")
	}
	return usr
}

func UpdateUserProfile(user structures.User, conn *sql.DB) (responses.UserListResponse, error) {
	var resp responses.UserListResponse

	stmt, err := conn.Prepare(queries.UPDATE_USER)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return resp, fmt.Errorf("Cannot Update User: %v", user.Id)
	}
	_, err = stmt.Exec(&user.Ig, &user.Twitter, &user.TikTok, &user.AvatarUrl,
		&user.ImgOne, &user.ImgTwo, &user.ImgThree, &user.Description, &user.Id)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	resp.Data = append(resp.Data, user)
	resp.Info.Count = len(resp.Data)
	resp.Info.Total = len(resp.Data)
	resp.Info.Next = false
	resp.Info.Page = 0
	return resp, nil
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

func GetUser(userId string, conn *sql.DB) (responses.UserListResponse, error) {
	var resp responses.UserListResponse
	uId, err := uuidtransform.StringToUuidTransform(userId)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	stmt, err := conn.Prepare(queries.GET_USER)
	if err != nil {
		return resp, err
	}
	defer stmt.Close()
	user := mapUserRow(stmt.QueryRow(&uId))
	resp.Data = append(resp.Data, user)
	resp.Info = structures.Info{Total: len(resp.Data), Count: len(resp.Data), Next: false, Page: 0}

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
	rows, err := stmt.Query(uId, time.Now().UTC(), util.CalcQueryStart(page), constants.LIMIT)
	if err != nil {
		log.Println(err)
		return resp, fmt.Errorf(`Unable to query member's event.`)
	}
	events, err := eventservices.MapMultiLineRows(rows)
	if err != nil {
		return resp, fmt.Errorf(`unable to map events`)
	}
	err = count.QueryRow(userId, time.Now().UTC()).Scan(&total)
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
	err = stmt.QueryRow(userId, time.Now().UTC()).Scan(&count)
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
	err = stmt.QueryRow(eventId, time.Now().UTC()).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
