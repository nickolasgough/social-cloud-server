package model

import (
	"time"
)

type Post struct {
	Username    string    `json:"username"`
	Displayname string    `json:"displayname"`
	Post        string    `json:"post"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE post (
	username VARCHAR(250) NOT NULL,
	post TEXT NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime),
	FOREIGN KEY (username) references profile (username)
);
`

const ModelDropQuery = `
DROP TABLE post;
`
