package model

import (
	"time"
)

type Post struct {
	Username string    `json:"username"`
	Avatar   Avatar    `json:"avatar"`
	Post     string    `json:"post"`
	Imageurl string    `json:"imageurl"`
	Likes    int       `json:"likes"`
	Dislikes int       `json:"dislikes"`
	Liked    bool      `json:"liked"`
	Disliked bool      `json:"disliked"`
	Datetime time.Time `json:"datetime"`
}

type Avatar struct {
	Displayname string `json:"displayname"`
	Imageurl    string `json:"imageurl"`
}

const ModelCreateQuery = `
CREATE TABLE post (
	username VARCHAR(250) NOT NULL,
	post TEXT NOT NULL,
	imageurl TEXT,
	likes INTEGER,
	dislikes INTEGER,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime),
	FOREIGN KEY (username) references profile (username)
);

CREATE TABLE reaction (
	username VARCHAR(250) NOT NULL,
	posttime TIMESTAMP NOT NULL,
	connection VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,
	reaction VARCHAR(250) NOT NULL,

	PRIMARY KEY (username, posttime, connection, datetime, reaction),
	FOREIGN KEY (username) references profile (username),
	FOREIGN KEY (connection) references profile (username),
	FOREIGN KEY (username, posttime) references post (username, datetime)
);
`

const ModelDropQuery = `
DROP TABLE reaction;

DROP TABLE post;
`
