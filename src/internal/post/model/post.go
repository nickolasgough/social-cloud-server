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
`

const ModelDropQuery = `
DROP TABLE post;
`
