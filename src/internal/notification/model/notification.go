package model

import (
	"time"
)


type Notification struct {
	Username    string    `json:"username"`
	Type        string    `json:"type"`
	Sender      string    `json:"sender"`
	Displayname string    `json:"displayname"`
	Dismissed   bool      `json:"dismissed"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE notification (
	username VARCHAR(250) NOT NULL,
	type VARCHAR(100) NOT NULL,
	sender VARCHAR(250) NOT NULL,
	dismissed BOOLEAN NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime),
	FOREIGN KEY (username) references profile (username),
	FOREIGN KEY (sender) references profile (username)
);
`

const ModelDropQuery = `
DROP TABLE notification;
`
