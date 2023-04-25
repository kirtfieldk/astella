package queries

var GET_EVENT_BY_ID_AND_LOCATION_INFO string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city FROM events e LEFT JOIN locationInfo l ON l.id = e.location_id WHERE e.id = $1`

var GET_EVENT_BY_CITY_AND_LOCATION_INFO string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city FROM events e LEFT JOIN locationInfo l ON l.id = e.location_id WHERE l.city = $1`

var GET_MESSAGES_IN_EVENT string = `SELECT id, content, user_id, created, event_id, parent_id,  upvotes,
	pinned, latitude, longitude FROM messages WHERE event_id = $1 ORDER BY created DESC LIMIT 30`

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

var FIND_USER_IN_EVENT = `SELECT id  FROM members WHERE user_id = $1 AND event_id = $2`
var INSERT_MESSAGE_WITH_PARENT_ID = `Insert INTO messages (content,user_id,created,event_id,parent_id,upvotes, pinned,latitude,longitude) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9);`
var INSERT_MESSAGE_WITHOUT_PARENT_ID = `Insert INTO messages (content,user_id,created,event_id,upvotes, pinned,latitude,longitude) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);`
var QUERY_ALL_WHO_LIKE_MESSAGE = `SELECT u.id, u.username, u.created, u.ig, u.twitter, u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three 
		FROM likes Left JOIN users u on likes.user_id = u.id WHERE likes.message_id = $1;`
