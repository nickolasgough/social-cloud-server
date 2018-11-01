package model

import (
	"time"
)


type Notification struct {
	Username  string
	Type      string
	Sender    string
	Dismissed bool
	Datetime  time.Time
}

const ModelCreateQuery = `
CREATE TABLE notification (
	username VARCHAR(250) NOT NULL,
	type VARCHAR(100) NOT NULL,
	sender VARCHAR(250) NOT NULL,
	dismissed BOOLEAN NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime)
);
`

const ModelDropQuery = `
DROP TABLE notification;
`
