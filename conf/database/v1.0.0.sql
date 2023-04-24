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
    location_id uuid,
    CONSTRAINT fk_location
      FOREIGN KEY(location_id) 
	  REFERENCES locationInfo(id)
    
);



CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
    created TIMESTAMP NOT NULL, 
    username varchar(80) NOT NULL,
    ig varchar(120),
    twitter varchar(120),
    tiktok varchar(120)
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
        REFERENCES events(id)
    
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
        REFERENCES events(id)
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
