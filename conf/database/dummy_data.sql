INSERT INTO locationInfo (top_left_lat,
    top_left_lon,
    top_right_lat,
    top_right_lon,
    bottom_left_lat,
    bottom_left_lon,
    bottom_right_lat,
    bottom_right_lon,
    city) VALUES 
    (40.7128,	-74.0060,	40.7128,	-73.9970,	40.7038,	-74.0060,	40.7038,	-73.9970,	'New York City'),
    (34.0522,	-118.2437,	34.0522,	-118.2347,	34.0432,	-118.2437,	34.0432,	-118.2347,	'Los Angeles');


INSERT INTO events (event_name, created, description, public, code, location_id)
VALUES ('Event 1', NOW(), 'Description of Event 1', TRUE, 'CODE1', (select id from locationInfo LIMIT 1)),
('Event 2', NOW(), 'Description of Event 2', FALSE, 'CODE2', (select id from locationInfo LIMIT 1)),
('Event 3', NOW(), 'Description of Event 3', TRUE, 'CODE3', (select id from locationInfo LIMIT 1));

-- Dummy data for users table
INSERT INTO users (username, created, ig, twitter, tiktok)
VALUES ('user1', NOW(), 'user1_ig', 'user1_twitter', 'user1_tiktok'),
('user2', NOW(), 'user2_ig', 'user2_twitter', 'user2_tiktok'),
('user3', NOW(), 'user3_ig', 'user3_twitter', 'user3_tiktok');

-- Dummy data for admins table
INSERT INTO admins (user_id, event_id, created)
VALUES ((SELECT id FROM users WHERE username = 'user1'), (SELECT id FROM events WHERE event_name = 'Event 1'), NOW()),
((SELECT id FROM users WHERE username = 'user2'), (SELECT id FROM events WHERE event_name = 'Event 2'),NOW()),
((SELECT id FROM users WHERE username = 'user3'), (SELECT id FROM events WHERE event_name = 'Event 3'),NOW());

-- Dummy data for members table
INSERT INTO members (user_id, event_id, created)
VALUES ((SELECT id FROM users WHERE username = 'user1'), (SELECT id FROM events WHERE event_name = 'Event 1'), NOW()),
((SELECT id FROM users WHERE username = 'user2'), (SELECT id FROM events WHERE event_name = 'Event 2'),NOW()),
((SELECT id FROM users WHERE username = 'user3'), (SELECT id FROM events WHERE event_name = 'Event 3'), NOW());

-- Dummy data for messages table
INSERT INTO messages (content, user_id, created, event_id, parent_message, upvotes, location_id)
VALUES ('Message 1 for Event 1', POINT(1,1), (SELECT id FROM users WHERE username = 'user1'), NOW(), (SELECT id FROM events WHERE event_name = 'Event 1'), NULL, 0, (select id from locationInfo LIMIT 1)),
('Message 2 for Event 1', POINT(2,2), (SELECT id FROM users WHERE username = 'user2'), NOW(), (SELECT id FROM events WHERE event_name = 'Event 1'), NULL, 0, (select id from locationInfo LIMIT 1)),
('Message 3 for Event 2', POINT(3,3), (SELECT id FROM users WHERE username = 'user3'), NOW(), (SELECT id FROM events WHERE event_name = 'Event 2'), NULL, 0, (select id from locationInfo LIMIT 1)),
('Message 4 for Event 2', POINT(4,4), (SELECT id FROM users WHERE username = 'user1'), NOW(), (SELECT id FROM events WHERE event_name = 'Event 2'), NULL, 0, (select id from locationInfo LIMIT 1)),
('Reply to Message 1', POINT(1,1), (SELECT id FROM users WHERE username = 'user2'), NOW(), (SELECT id FROM events WHERE event_name = 'Event 1'), (SELECT id FROM messages WHERE content = 'Message 1 for Event 1'), 0, (select id from locationInfo LIMIT 1));