package queries

var GET_EVENT_BY_ID_AND_LOCATION_INFO string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city, e.expired, e.end_time FROM events e LEFT JOIN locationInfo l ON l.id = e.location_id WHERE e.id = $1
	and e.end_time >= $2;`

var GET_EVENT_BY_CITY_AND_LOCATION_INFO string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city, e.expired, e.end_time FROM events e LEFT JOIN locationInfo l ON l.id = e.location_id WHERE l.city = $1 
	AND e.end_time >= $2 OFFSET $3 LIMIT $4;`

var GET_EVENT_BY_CITY_AND_LOCATION_INFO_COUNT string = `Select COUNT(e.id) FROM events e LEFT JOIN 
	locationInfo l ON l.id = e.location_id WHERE l.city = $1 AND e.end_time >= $2;`

var GET_EVENTS_LOCATION_INFO_USER_IN string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city, e.expired, e.end_time FROM members m LEFT JOIN events e on m.event_id = e.id
	LEFT JOIN locationInfo l ON l.id = e.location_id WHERE m.user_id = $1 AND e.end_time >= $2 OFFSET $3 LIMIT $4;`

var GET_EVENTS_MEMBER_OF_COUNT string = `Select COUNT(members.id) FROM members LEFT JOIN events e 
	ON e.id = members.event_id WHERE user_id = $1 AND e.end_time >= $2;`

var GET_EVENT_USERS string = `Select u.id, u.username, u.description, u.created, u.ig, u.twitter, u.tiktok, u.avatar_url, 
	u.img_one, u.img_two, u.img_three FROM members m LEFT JOIN users u on m.user_id = u.id
	WHERE m.event_id = $1 OFFSET $2 LIMIT $3;`

var GET_EVENTS_LOCATION_INFO_USER_IN_COUNT string = `Select COUNT(e.id) FROM members m LEFT JOIN events e on m.event_id = e.id
	WHERE m.user_id = $1 AND e.end_time >= $2;`
var GET_EVENTS_MEMBERS_COUNT string = `Select COUNT(m.id) FROM members m LEFT JOIN events e on m.event_id = e.id
	WHERE e.id = $1 AND e.end_time >= $2;`

var GET_MESSAGES_IN_EVENT string = `SELECT m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
	m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
	u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three FROM messages m LEFT JOIN users u ON m.user_id = u.id 
	WHERE event_id = $1 ORDER BY m.created DESC OFFSET $2 LIMIT $3;`

var GET_MESSAGES_IN_EVENT_COUNT string = `SELECT COUNT(id) FROM messages WHERE event_id = $1;`

var GET_LOCATION_FOR_EVENT string = `SELECT  
	locationInfo.id, 
	top_left_lat,
    top_left_lon,
    top_right_lat,
    top_right_lon,
    bottom_left_lat,
    bottom_left_lon,
    bottom_right_lat,
    bottom_right_lon
	FROM locationInfo LEFT JOIN events ON 
	events.location_id = locationInfo.id WHERE events.id = $1`

var FIND_IF_USER_IN_EVENT string = `SELECT id  FROM members WHERE user_id = $1 AND event_id = $2`
var INSERT_MESSAGE_WITH_PARENT_ID string = `Insert INTO messages (content,user_id,created,event_id,parent_id,upvotes, pinned,latitude,longitude) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9);`
var INSERT_MESSAGE_WITHOUT_PARENT_ID string = `Insert INTO messages (content,user_id,created,event_id,upvotes, pinned,latitude,longitude) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);`
var QUERY_ALL_WHO_LIKE_MESSAGE string = `SELECT u.id, u.username, u.description, u.created, u.ig, u.twitter, u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three 
		FROM likes Left JOIN users u on likes.user_id = u.id WHERE likes.message_id = $1 OFFSET $2 LIMIT $3;`
var QUERY_ALL_WHO_LIKE_MESSAGE_COUNT string = `SELECT COUNT(id) FROM likes WHERE likes.message_id = $1;`
var INSERT_EVENT_INTO_DB = `INSERT INTO events (event_name, created, description, public, code, location_id, duration, end_time, expired) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`
var INSERT_LOCATION_INTO_DB_RETURN_ID string = `INSERT INTO locationInfo (top_left_lat, top_left_lon, top_right_lat, 
	top_right_lon, bottom_left_lat, bottom_left_lon, bottom_right_lat, bottom_right_lon, city) VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id;`

var EXPIRE_EVENT string = `UPDATE events SET expired = true WHERE id = $1`
var INSERT_UPVOTE string = `INSERT INTO likes (user_id, message_id, created) VALUES ($1,$2, $3)`
var DELETE_UPVOTE string = `DELETE FROM likes WHERE user_id = $1 AND message_id = $2;`

var UPDATE_MESSAGE_LIKE_INC string = `UPDATE messages SET upvotes = upvotes + 1 WHERE id = $1`
var UPDATE_MESSAGE_LIKE_DEC string = `UPDATE messages SET upvotes = upvotes - 1 WHERE id = $1`
