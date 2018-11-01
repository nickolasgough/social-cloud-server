package model

import (
	"time"
)

type Post struct {
	Username string
	Post     string
	Datetime time.Time
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
