package queries

const GET_EVENT_BY_ID_AND_LOCATION_INFO string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city, e.user_id, e.end_time FROM events e LEFT JOIN locationInfo l ON l.id = e.location_id WHERE e.id = $1
	and e.end_time >= $2;`

const GET_EVENT_BY_CITY_AND_LOCATION_INFO string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city, e.user_id, e.end_time FROM events e LEFT JOIN locationInfo l ON l.id = e.location_id WHERE l.city = $1 
	AND e.end_time >= $2 OFFSET $3 LIMIT $4;`

const GET_EVENT_BY_CITY_AND_LOCATION_INFO_COUNT string = `Select COUNT(e.id) FROM events e LEFT JOIN 
	locationInfo l ON l.id = e.location_id WHERE l.city = $1 AND e.end_time >= $2;`

const GET_EVENTS_LOCATION_INFO_USER_IN string = `Select e.id, e.event_name, e.created, e.description, e.public, e.code, 
	l.id, l.top_left_lat, l.top_left_lon, l.top_right_lat, l.top_right_lon, l.bottom_right_lat,l.bottom_right_lon,
	l.bottom_left_lat, l.bottom_left_lon, l.city, e.user_id, e.end_time FROM members m LEFT JOIN events e on m.event_id = e.id
	LEFT JOIN locationInfo l ON l.id = e.location_id WHERE m.user_id = $1 AND e.end_time >= $2 OFFSET $3 LIMIT $4;`

const GET_EVENTS_MEMBER_OF_COUNT string = `Select COUNT(members.id) FROM members LEFT JOIN events e 
	ON e.id = members.event_id WHERE members.user_id = $1 AND e.end_time >= $2;`

const GET_EVENT_USERS string = `Select u.id, u.username, u.description, u.created, u.ig, u.twitter, u.tiktok, u.avatar_url, 
	u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube FROM members m LEFT JOIN users u on m.user_id = u.id
	WHERE m.event_id = $1 OFFSET $2 LIMIT $3;`

const GET_EVENTS_LOCATION_INFO_USER_IN_COUNT string = `Select COUNT(e.id) FROM members m LEFT JOIN events e on m.event_id = e.id
	WHERE m.user_id = $1 AND e.end_time >= $2;`
const GET_EVENTS_MEMBERS_COUNT string = `Select COUNT(m.id) FROM members m LEFT JOIN events e on m.event_id = e.id
	WHERE e.id = $1 AND e.end_time >= $2;`

const GET_MESSAGES_IN_EVENT string = `SELECT m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
	m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
	u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube, COUNT(m.id)  FROM messages m LEFT JOIN users u ON m.user_id = u.id 
	LEFT JOIN messages m2 on m.id = m2.parent_id
	WHERE m.event_id = $1 AND m.parent_id IS NULL GROUP BY m.id, u.id ORDER BY m.created DESC OFFSET $2 LIMIT $3;`

const GET_MESSAGES string = `SELECT m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
	m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
	u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube, COUNT(m.id) FROM messages m LEFT JOIN users u ON m.user_id = u.id 
	LEFT JOIN users u ON m.user_id = u.id 
	LEFT JOIN messages m2 on m.id = m2.parent_id
	WHERE m.id = $1 GROUP BY m.id, u.id;`

const GET_MESSAGES_IN_EVENT_COUNT string = `SELECT COUNT(id) FROM messages WHERE event_id = $1 AND parent_id IS NULL;`

const GET_LOCATION_FOR_EVENT string = `SELECT  
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

const FIND_IF_USER_IN_EVENT string = `SELECT id FROM public.members WHERE user_id = $1 AND event_id = $2`
const INSERT_MESSAGE_WITH_PARENT_ID string = `Insert INTO messages (content,user_id,created,event_id,parent_id,upvotes, pinned,latitude,longitude) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9);`
const INSERT_MESSAGE_WITHOUT_PARENT_ID string = `Insert INTO messages (content,user_id,created,event_id,upvotes, pinned,latitude,longitude) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);`
const QUERY_ALL_WHO_LIKE_MESSAGE string = `SELECT u.id, u.username, u.description, u.created, u.ig, u.twitter, u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three,
		u.snapchat,u.youtube
		FROM likes Left JOIN users u on likes.user_id = u.id WHERE likes.message_id = $1 OFFSET $2 LIMIT $3;`
const QUERY_ALL_WHO_LIKE_MESSAGE_COUNT string = `SELECT COUNT(id) FROM likes WHERE likes.message_id = $1;`

const INSERT_EVENT_INTO_DB = `INSERT INTO events (event_name, created, description, public, code, location_id, duration, end_time, user_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id;`

const INSERT_LOCATION_INTO_DB_RETURN_ID string = `INSERT INTO locationInfo (top_left_lat, top_left_lon, top_right_lat, 
	top_right_lon, bottom_left_lat, bottom_left_lon, bottom_right_lat, bottom_right_lon, city) VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id;`

const EXPIRE_EVENT string = `UPDATE events SET expired = true WHERE id = $1`
const INSERT_UPVOTE string = `INSERT INTO likes (user_id, message_id, created) VALUES ($1,$2, $3)`
const DELETE_UPVOTE string = `DELETE FROM likes WHERE user_id = $1 AND message_id = $2;`
const INSERT_PINNED string = `INSERT INTO pinned (user_id, message_id, created) VALUES ($1,$2, $3)`
const DELETE_PINNED string = `DELETE FROM pinned WHERE user_id = $1 AND message_id = $2;`

const UPDATE_MESSAGE_LIKE_INC string = `UPDATE messages m SET upvotes = upvotes + 1 FROM users u
 WHERE m.id = $1 and m.user_id = u.id
	RETURNING m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
	m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
	u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube;`

const UPDATE_MESSAGE_LIKE_DEC string = `UPDATE messages m SET upvotes = upvotes - 1 FROM users u
 WHERE m.id = $1 and m.user_id = u.id
	RETURNING m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
	m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
	u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube;`

const GET_USER_PIN_MSG_IN_EVENT string = `SELECT 
		m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
		m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
		u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube
	FROM pinned p LEFT JOIN messages m on p.message_id = m.id LEFT JOIN users u
	on u.id = m.user_id WHERE p.user_id = $1 and m.event_id = $2 ORDER BY p.created DESC OFFSET $3 LIMIT $4;`

const GET_USER_PIN_MSG_IN_EVENT_COUNT string = `SELECT 
	COUNT(p.id)
	FROM pinned p LEFT JOIN messages m on p.message_id = m.id LEFT JOIN users u
	on u.id = m.user_id WHERE p.user_id = $1 and m.event_id = $2;`

const GET_USER_SINGLE_PIN_MSG_IN_EVENT string = `SELECT 
	COUNT(p.id) FROM pinned p WHERE p.user_id = $1 and p.message_id = $2;`

const GET_USER_SINGLE_LIKE_MSG_IN_EVENT string = `SELECT 
	COUNT(l.id) FROM likes l WHERE l.user_id = $1 and l.message_id = $2;`

const GET_MESSAGE_THREAD string = `SELECT m.id, m.content, m.created, m.event_id, m.parent_id,  m.upvotes,
	m.pinned, m.latitude, m.longitude, u.id, u.username, u.description, u.created, u.ig, u.twitter, 
	u.tiktok, u.avatar_url, u.img_one, u.img_two, u.img_three, u.snapchat, u.youtube, COUNT(m.id) FROM messages m LEFT JOIN users u ON m.user_id = u.id 
	LEFT JOIN messages m2 on m.id = m2.parent_id
	WHERE m.parent_id = $1 GROUP BY m.id, u.id OFFSET $3 LIMIT $2; `

const GET_MESSAGE_THREAD_COUNT string = `SELECT COUNT(id) FROM messages WHERE parent_id = $1; `

const UPDATE_USER string = `UPDATE users SET ig = $1, twitter = $2, tiktok = $3, avatar_url = $4, 
		img_one = $5, img_two = $6, img_three = $7, description = $8 WHERE id = $9`

const GET_USER string = `SELECT id, username, description, created, ig, twitter, tiktok, avatar_url, img_one, img_two, img_three,
	snapchat, youtube FROM users where id = $1`

const ADD_USER_TO_EVENT string = `INSERT INTO members (user_id, event_id, created) values ($1, $2, $3);`
