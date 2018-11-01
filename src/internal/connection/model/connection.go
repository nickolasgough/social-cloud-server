package model

import (
	"time"
)


type Connection struct {
	Username string
	Sender   string
	Datetime time.Time
}

const ModelCreateQuery = `
CREATE TABLE connection (
	username VARCHAR(250) NOT NULL,
	sender VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime)
);
`

const ModelDropQuery = `
DROP TABLE connection;
`
