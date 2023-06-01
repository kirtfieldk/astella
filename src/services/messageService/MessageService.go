package messageservice

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/api/responses"
	"github.com/kirtfieldk/astella/src/constants"
	"github.com/kirtfieldk/astella/src/constants/queries"
	locationservice "github.com/kirtfieldk/astella/src/services/locationService"
	userservice "github.com/kirtfieldk/astella/src/services/userService"
	"github.com/kirtfieldk/astella/src/structures"
	"github.com/kirtfieldk/astella/src/util"
	uuidtransform "github.com/kirtfieldk/astella/src/util/uuidTransform"
)

func PostMessage(msg structures.MessageRequestBody, conn *sql.DB) (bool, error) {
	eventId, userId, err := uuidtransform.ParseTwoIds(msg.EventId, msg.UserId)
	if err != nil {
		log.Printf("Failed to be UUID for Event: " + msg.EventId)
		return false, err
	}
	if isRequestInArea(eventId, msg.Latitude, msg.Longitude, conn) && isUserInEvent(userId, eventId, conn) {
		return insertMessageIntoDB(msg, userId, eventId, conn), nil
	}
	log.Printf("User %v is not in event %v\n", msg.UserId, msg.EventId)
	return false, nil

}

func GetMessagesInEvent(eventId string, userId string, pagination int, conn *sql.DB) (responses.MessageListResponse, error) {
	var messages responses.MessageListResponse
	eId, uId, err := uuidtransform.ParseTwoIds(eventId, userId)
	if err != nil {
		return messages, err
	}
	if isUserInEvent(uId, eId, conn) {
		messages, err = getMessages(eId, pagination, conn)
		if err != nil {
			return messages, err
		}
		populateMsgPinnedByUser(messages.Data, uId, conn)
		populateMsgLikedByUser(messages.Data, uId, conn)
		return messages, nil
	}
	return messages, fmt.Errorf("User is not in Event")

}

func UpvoteMessage(messageId string, userId string, eventId string, cords structures.Point, conn *sql.DB) (responses.MessageListResponse, error) {
	var resp responses.MessageListResponse
	eId, uId, mId, err := uuidtransform.ParseThreeIds(eventId, userId, messageId)
	if err != nil {
		return resp, err
	}
	if isRequestInArea(eId, cords.Latitude, cords.Longitude, conn) && isUserInEvent(uId, eId, conn) {
		msg, err := upVoteMessage(uId, mId, conn)
		if err != nil {
			return resp, err
		}
		resp.Data = append(resp.Data, msg)
		populateMsgPinnedByUser(resp.Data, uId, conn)
		populateMsgLikedByUser(resp.Data, uId, conn)
		resp.Info = structures.Info{Total: 1, Count: 1, Next: false}
		return resp, nil
	}
	return resp, fmt.Errorf("User is not in Message Event: ")
}

func DownVoteMessage(messageId string, userId string, eventId string, cords structures.Point, conn *sql.DB) (responses.MessageListResponse, error) {
	var resp responses.MessageListResponse
	eId, uId, mId, err := uuidtransform.ParseThreeIds(eventId, userId, messageId)
	if err != nil {
		return resp, err
	}
	if isRequestInArea(eId, cords.Latitude, cords.Longitude, conn) && isUserInEvent(uId, eId, conn) {
		event, err := downVoteMessage(uId, mId, conn)
		if err != nil {
			return resp, err
		}
		resp.Data = append(resp.Data, event)
		resp.Info = structures.Info{Total: 1, Count: 1, Next: false}
		populateMsgPinnedByUser(resp.Data, uId, conn)
		populateMsgLikedByUser(resp.Data, uId, conn)
		return resp, nil
	}
	return resp, fmt.Errorf("User is not in Message Event: ")
}

func GetMessageThread(messageId string, userId string, eventId string, page int, conn *sql.DB) (responses.MessageListResponse, error) {
	var resp responses.MessageListResponse
	var ct int
	mId, uId, eId, err := uuidtransform.ParseThreeIds(messageId, userId, eventId)
	if err != nil {
		return resp, err
	}
	if isUserInEvent(uId, eId, conn) {
		log.Println(mId)
		stmt, err := conn.Prepare(queries.GET_MESSAGE_THREAD)
		defer stmt.Close()
		if err != nil {
			log.Println(err)
			return resp, err
		}
		count, err := conn.Prepare(queries.GET_MESSAGE_THREAD_COUNT)
		defer count.Close()
		if err != nil {
			return resp, err
		}
		rows, err := stmt.Query(mId, constants.LIMIT, util.CalcQueryStart(page))
		if err != nil {
			log.Println(err)
			return resp, err
		}
		if err = count.QueryRow(&mId).Scan(&ct); err != nil {
			return resp, err
		}

		messages := mapRowsToMessages(rows)
		populateMsgLikedByUser(messages, uId, conn)
		populateMsgPinnedByUser(messages, uId, conn)
		resp.Data = messages
		resp.Info.Count = len(messages)
		resp.Info.Total = ct
		resp.Info.Page = page
		resp.Info.Next = page*constants.LIMIT < ct && len(messages) != ct
		return resp, nil
	}
	return resp, fmt.Errorf("User not a member of event")
}

func GetUserUpvotes(messageId string, userId string, eventId string, pagination int, conn *sql.DB) (responses.UserListResponse, error) {
	var userResponse responses.UserListResponse
	eId, uId, mId, err := uuidtransform.ParseThreeIds(eventId, userId, messageId)
	if err != nil {
		return userResponse, err
	}
	if isUserInEvent(uId, eId, conn) {
		users, err := getUsersUpvoteMessage(mId, pagination, conn)
		if err != nil {
			return userResponse, err
		}
		total, err := totalAmountOfUpvotes(mId, conn)
		if err != nil {
			return userResponse, err
		}
		info := structures.Info{
			Count: len(users),
			Page:  pagination,
			Total: total,
			Next:  pagination*constants.LIMIT < total && len(users) != total,
		}

		userResponse.Info = info
		userResponse.Data = users
		return userResponse, nil
	}
	return userResponse, fmt.Errorf("User is not in Event")
}

func GetUsersPinnedMessagesInEvent(userId string, eventId string, page int, conn *sql.DB) (responses.MessageListResponse, error) {
	var resp responses.MessageListResponse
	var total int
	eId, uId, err := uuidtransform.ParseTwoIds(eventId, userId)
	if err != nil {
		log.Printf("Failed to be UUID for User: " + userId)
		return resp, err
	}
	stmt, err := conn.Prepare(queries.GET_USER_PIN_MSG_IN_EVENT)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return resp, err
	}
	rows, err := stmt.Query(uId, eId, util.CalcQueryStart(page), constants.LIMIT)
	if err != nil {
		return resp, err
	}
	messages := mapRowsToMessages(rows)
	countStmt, err := conn.Prepare(queries.GET_USER_PIN_MSG_IN_EVENT_COUNT)
	defer countStmt.Close()
	err = countStmt.QueryRow(uId, eId).Scan(&total)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	populateMsgPinnedByUser(messages, uId, conn)
	populateMsgLikedByUser(messages, uId, conn)
	resp.Data = messages
	resp.Info.Total = total
	resp.Info.Count = len(messages)
	resp.Info.Next = page*constants.LIMIT < total && len(messages) != total
	return resp, nil
}

func PinnMessage(messageId string, userId string, eventId string, conn *sql.DB) (responses.MessageListResponse, error) {
	var resp responses.MessageListResponse
	eId, uId, mId, err := uuidtransform.ParseThreeIds(eventId, userId, messageId)
	if err != nil {
		return resp, err
	}
	if isUserInEvent(uId, eId, conn) {
		stmt, err := conn.Prepare(queries.INSERT_PINNED)
		defer stmt.Close()
		if err != nil {
			log.Println(err)
			return resp, err
		}
		_, err = stmt.Exec(&uId, &mId, time.Now().UTC())
		if err != nil {
			log.Println(err)
			return resp, err
		}
		messages, err := getMessage(mId, conn)
		if err != nil {
			log.Println(err)
			return resp, nil
		}
		populateMsgLikedByUser(messages, uId, conn)
		populateMsgPinnedByUser(messages, uId, conn)
		resp.Data = messages
		resp.Info.Count = 1
		resp.Info.Total = 1
		resp.Info.Next = false
	}
	return resp, nil
}

func UnpinnMessage(messageId string, userId string, eventId string, conn *sql.DB) (responses.MessageListResponse, error) {
	var resp responses.MessageListResponse
	eId, uId, mId, err := uuidtransform.ParseThreeIds(eventId, userId, messageId)
	if err != nil {
		return resp, err
	}
	if isUserInEvent(uId, eId, conn) {
		stmt, err := conn.Prepare(queries.DELETE_PINNED)
		defer stmt.Close()
		if err != nil {
			log.Println(err)
			return resp, err
		}
		_, err = stmt.Exec(&uId, &mId)
		if err != nil {
			log.Println(err)
			return resp, err
		}
		messages, err := getMessage(mId, conn)
		if err != nil {
			log.Println(err)
			return resp, nil
		}
		populateMsgLikedByUser(messages, uId, conn)
		populateMsgPinnedByUser(messages, uId, conn)
		resp.Data = messages
		resp.Info.Count = 1
		resp.Info.Total = 1
		resp.Info.Next = false
	}
	return resp, nil
}

func getMessage(messageId uuid.UUID, conn *sql.DB) ([]structures.Message, error) {
	var msg []structures.Message
	stmt, err := conn.Prepare(queries.GET_MESSAGES)
	if err != nil {
		return msg, err
	}
	rows, err := stmt.Query(messageId)
	if err != nil {
		return msg, err
	}
	return mapRowsToMessages(rows), nil
}

func getUsersUpvoteMessage(messageId uuid.UUID, pagination int, conn *sql.DB) ([]structures.User, error) {
	rows, err := conn.Query(queries.QUERY_ALL_WHO_LIKE_MESSAGE, messageId, util.CalcQueryStart(pagination), constants.LIMIT)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return userservice.MapUserRows(rows), nil
}

func totalAmountOfUpvotes(messageId uuid.UUID, conn *sql.DB) (int, error) {
	var count int
	if err := conn.QueryRow(queries.QUERY_ALL_WHO_LIKE_MESSAGE_COUNT, messageId).Scan(&count); err != nil {
		log.Println(err)
		return count, fmt.Errorf("Unable to count likes for %v", messageId)
	}
	return count, nil
}

func getErrorMessageForNoMessages(err error, id string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("Messages For EventId %s: no messages", id)
	}
	return fmt.Errorf("MessagesEventById %s: %v", id, err)
}

func mapRowsToMessages(rows *sql.Rows) []structures.Message {
	var messages []structures.Message = make([]structures.Message, 0)
	for rows.Next() {
		var msg structures.Message
		var usr structures.User
		var parentId sql.NullString
		if err := rows.Scan(
			&msg.Id, &msg.Content, &msg.Created, &msg.EventId, &parentId, &msg.Upvotes, &msg.Pinned, &msg.Latitude, &msg.Longitude,
			&usr.Id, &usr.Username, &usr.Description, &usr.Created, &usr.Ig, &usr.Twitter, &usr.TikTok, &usr.AvatarUrl,
			&usr.ImgOne, &usr.ImgTwo, &usr.ImgThree, &usr.Snapchat, &usr.Youtube, &msg.Replies); err != nil {
			log.Println(err)
			log.Println("Issue mapping DB row")
		}
		if &parentId != nil {
			msg.ParentId = parentId.String
		}
		msg.User = usr
		messages = append(messages, msg)
	}

	return messages
}

func isRequestInArea(id uuid.UUID, lat float32, long float32, conn *sql.DB) bool {
	location := locationservice.GetEventLocation(id, conn)
	topLeftPt := locationservice.CreatePoint(location.TopLeftLat, location.TopLeftLon)
	topRightPt := locationservice.CreatePoint(location.TopRightLat, location.TopRightLon)
	bottomLeftPt := locationservice.CreatePoint(location.BottomLeftLat, location.BottomLeftLon)
	bottomRightPt := locationservice.CreatePoint(location.BottomLeftLat, location.BottomRightLon)
	return locationservice.CheckPointInArea(lat, long, topLeftPt, topRightPt, bottomLeftPt, bottomRightPt)
}

//9. Instead goea from router to platform
//Check how longtoken is valid in cache
// Enrich and capture - Discuss capture in here
//Customer Access -> Deployed in CloudFront
//Add more text for UI Data Flow ->
//When a message hits summarists, How does it make its way through Apollo applications
//Add ALB in in UI Data Flow

func isUserInEvent(userId uuid.UUID, eventId uuid.UUID, conn *sql.DB) bool {
	var id uuid.UUID
	stmt, err := conn.Prepare(queries.FIND_IF_USER_IN_EVENT)
	defer stmt.Close()
	if err != nil {
		log.Println("Touble preparing statement for isUserInEvent")
	}
	err = stmt.QueryRow(userId, eventId).Scan(&id)
	if err == sql.ErrNoRows {
		log.Println("No rpos")
		return false
	}
	return true
}

func insertMessageIntoDB(msg structures.MessageRequestBody, userId uuid.UUID, eventId uuid.UUID, conn *sql.DB) bool {
	if msg.ParentId != "" {
		_, err := conn.Exec(queries.INSERT_MESSAGE_WITH_PARENT_ID, &msg.Content, userId, time.Now().UTC(), eventId, &msg.ParentId, &msg.Upvotes, &msg.Pinned,
			&msg.Latitude, &msg.Longitude)
		if err != nil {
			log.Println(err)
			return false
		}
	} else {
		_, err := conn.Exec(queries.INSERT_MESSAGE_WITHOUT_PARENT_ID, &msg.Content, userId, time.Now().UTC(), eventId, &msg.Upvotes, &msg.Pinned,
			&msg.Latitude, &msg.Longitude)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}

func getMessages(eventId uuid.UUID, page int, conn *sql.DB) (responses.MessageListResponse, error) {
	var response responses.MessageListResponse
	var messages []structures.Message
	var total int
	rows, err := conn.Query(queries.GET_MESSAGES_IN_EVENT, eventId, util.CalcQueryStart(page), constants.LIMIT)
	defer rows.Close()
	if err != nil {
		return response, err
	}
	messages = mapRowsToMessages(rows)
	if err = conn.QueryRow(queries.GET_MESSAGES_IN_EVENT_COUNT, eventId).Scan(&total); err != nil {
		return response, err
	}
	response.Data = messages
	response.Info.Total = total
	response.Info.Count = len(messages)
	response.Info.Next = page*constants.LIMIT < total && len(messages) != total

	return response, nil

}

func upVoteMessage(userId uuid.UUID, messageId uuid.UUID, conn *sql.DB) (structures.Message, error) {
	var msg structures.Message
	var usr structures.User
	var parentId sql.NullString
	tx, err := conn.Begin()
	stmt, err := tx.Prepare(queries.INSERT_UPVOTE)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return msg, fmt.Errorf(`Issue Preparing Upvote Insert Stmt`)
	}
	if _, err = stmt.Exec(&userId, &messageId, time.Now().UTC()); err != nil {
		log.Println(err)
		return msg, fmt.Errorf(constants.UNABLE_TO_UPVOTE_FOR_USER, userId)
	}
	updateMsgStmt, err := tx.Prepare(queries.UPDATE_MESSAGE_LIKE_INC)
	defer updateMsgStmt.Close()
	if err != nil {
		log.Panicln(err)
		return msg, fmt.Errorf(`Issue Preparing Upvote Stmt`)
	}
	if err = updateMsgStmt.QueryRow(&messageId).Scan(&msg.Id, &msg.Content, &msg.Created, &msg.EventId, &parentId, &msg.Upvotes, &msg.Pinned, &msg.Latitude, &msg.Longitude,
		&usr.Id, &usr.Username, &usr.Description, &usr.Created, &usr.Ig, &usr.Twitter, &usr.TikTok, &usr.AvatarUrl,
		&usr.ImgOne, &usr.ImgTwo, &usr.ImgThree, &usr.Snapchat, &usr.Youtube); err != nil {
		log.Println(err)
		return msg, fmt.Errorf(constants.UNABLE_TO_UPVOTE_FOR_USER, userId)
	}
	if &parentId != nil {
		msg.ParentId = parentId.String
	}
	tx.Commit()
	msg.User = usr
	return msg, nil
}

func downVoteMessage(userId uuid.UUID, messageId uuid.UUID, conn *sql.DB) (structures.Message, error) {
	var msg structures.Message
	var usr structures.User
	var parentId sql.NullString

	tx, err := conn.Begin()
	stmt, err := tx.Prepare(queries.DELETE_UPVOTE)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return msg, fmt.Errorf(`Issue Preparing Upvote Delete Stmt`)
	}
	if _, err = stmt.Exec(&userId, &messageId); err != nil {
		log.Println(err)
		return msg, fmt.Errorf(constants.UNABLE_TO_UPVOTE_FOR_USER, userId)
	}
	updateMsgStmt, err := tx.Prepare(queries.UPDATE_MESSAGE_LIKE_DEC)
	defer updateMsgStmt.Close()
	if err != nil {
		log.Println(err)
		return msg, fmt.Errorf(`Issue Preparing DownVote Stmt`)
	}
	if err = updateMsgStmt.QueryRow(&messageId).Scan(&msg.Id, &msg.Content, &msg.Created, &msg.EventId, &parentId, &msg.Upvotes, &msg.Pinned, &msg.Latitude, &msg.Longitude,
		&usr.Id, &usr.Username, &usr.Description, &usr.Created, &usr.Ig, &usr.Twitter, &usr.TikTok, &usr.AvatarUrl,
		&usr.ImgOne, &usr.ImgTwo, &usr.ImgThree, &usr.Snapchat, &usr.Youtube); err != nil {
		log.Println(err)
		return msg, fmt.Errorf(constants.UNABLE_TO_UPVOTE_FOR_USER, userId)
	}
	if &parentId != nil {
		msg.ParentId = parentId.String
	}
	tx.Commit()
	msg.User = usr
	return msg, nil
}

func populateMsgPinnedByUser(messages []structures.Message, userId uuid.UUID, conn *sql.DB) {
	stmt, err := conn.Prepare(queries.GET_USER_SINGLE_PIN_MSG_IN_EVENT)
	defer stmt.Close()
	if err != nil {
		return
	}
	for i := 0; i < len(messages); i++ {
		var ct int
		stmt.QueryRow(&userId, &messages[i].Id).Scan(&ct)
		if ct == 1 {
			messages[i].PinnedByUser = true
		} else {
			messages[i].PinnedByUser = false
		}
	}
}

func populateMsgLikedByUser(messages []structures.Message, userId uuid.UUID, conn *sql.DB) {
	stmt, err := conn.Prepare(queries.GET_USER_SINGLE_LIKE_MSG_IN_EVENT)
	defer stmt.Close()
	if err != nil {
		return
	}
	for i := 0; i < len(messages); i++ {
		var ct int
		stmt.QueryRow(&userId, &messages[i].Id).Scan(&ct)
		if ct == 1 {
			messages[i].UpvotedByUser = true
		} else {
			messages[i].UpvotedByUser = false
		}
	}

}
