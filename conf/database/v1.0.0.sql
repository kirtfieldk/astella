-- CREATE DATABASE astella;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS locationInfo(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    top_left_lat float,
    top_left_lon float,
    top_right_lat float,
    top_right_lon float,
    bottom_left_lat float,
    bottom_left_lon float,
    bottom_right_lat float,
    bottom_right_lon float,
    city varchar(100)
);

 CREATE TABLE IF NOT EXISTS events(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    event_name varchar(80) NOT NULL,
    created TIMESTAMP NOT NULL,
    description varchar(400),
    public BOOLEAN DEFAULT FALSE,
    code varchar(20),
    expired BOOLEAN DEFAULT FALSE,
    end_time TIMESTAMP NOT NULL, 
    duration float,
    location_id uuid,
    CONSTRAINT fk_location
      FOREIGN KEY(location_id) 
	  REFERENCES locationInfo(id)
    
);



CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    created TIMESTAMP NOT NULL, 
    username varchar(80) NOT NULL,
    description varchar(300),
    ig varchar(120),
    twitter varchar(120),
    tiktok varchar(120),
    avatar_url varchar(120),
    img_one varchar(120),
    img_two varchar(120),
    img_three varchar(120)
);

CREATE TABLE IF NOT EXISTS admins(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    created TIMESTAMP NOT NULL,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id),
    CONSTRAINT fk_event
        FOREIGN KEY(event_id) 
        REFERENCES events(id),
    UNIQUE (user_id, event_id)
    
);

CREATE TABLE IF NOT EXISTS members (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    created TIMESTAMP NOT NULL,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id),
    CONSTRAINT fk_event
        FOREIGN KEY(event_id) 
        REFERENCES events(id),
    UNIQUE (user_id, event_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    content text NOT NULL,
    user_id UUID NOT NULL,
    created TIMESTAMP NOT NULL,
    event_id UUID NOT NULL,
    parent_id UUID,
    upvotes integer DEFAULT 0,
    pinned BOOLEAN DEFAULT FALSE,
    latitude float,
    longitude float,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id),
    CONSTRAINT fk_event
        FOREIGN KEY(event_id) 
        REFERENCES events(id)

);

CREATE TABLE IF NOT EXISTS likes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (), 
    user_id UUID NOT NULL,
    message_id UUID NOT NULL,
    created TIMESTAMP,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id),
    CONSTRAINT fk_message
        FOREIGN KEY(message_id) 
        REFERENCES messages(id),
    UNIQUE (user_id, message_id)
);


-- Dummy data for users table
INSERT INTO users (username, created, ig, twitter, tiktok)
VALUES ('user1', NOW(), 'user1_ig', 'user1_twitter', 'user1_tiktok'),
('user2', NOW(), 'user2_ig', 'user2_twitter', 'user2_tiktok'),
('user3', NOW(), 'user3_ig', 'user3_twitter', 'user3_tiktok');