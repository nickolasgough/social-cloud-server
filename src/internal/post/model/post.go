package model

import (
	"time"
)

type Post struct {
	Email    string    `json:"email"`
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
	email VARCHAR(250) NOT NULL,
	post TEXT NOT NULL,
	imageurl TEXT,
	likes INTEGER,
	dislikes INTEGER,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (email, datetime),
	FOREIGN KEY (email) references profile (email)
);

CREATE TABLE reaction (
	email VARCHAR(250) NOT NULL,
	posttime TIMESTAMP NOT NULL,
	connection VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,
	reaction VARCHAR(250) NOT NULL,

	PRIMARY KEY (email, posttime, connection, datetime, reaction),
	FOREIGN KEY (email) references profile (email),
	FOREIGN KEY (connection) references profile (email),
	FOREIGN KEY (email, posttime) references post (email, datetime)
);
`

const ModelDropQuery = `
DROP TABLE reaction;

DROP TABLE post;
`
