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
	text TEXT NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime)
);
`

const ModelDropQuery = `
DROP TABLE post;
`
