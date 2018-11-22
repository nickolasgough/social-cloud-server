package model

import (
	"time"
)

type Notification struct {
	Email       string    `json:"email"`
	Type        string    `json:"type"`
	Sender      string    `json:"sender"`
	Displayname string    `json:"displayname"`
	Dismissed   bool      `json:"dismissed"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE notification (
	email VARCHAR(250) NOT NULL,
	type VARCHAR(100) NOT NULL,
	sender VARCHAR(250) NOT NULL,
	dismissed BOOLEAN NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (email, datetime),
	FOREIGN KEY (email) references profile (email),
	FOREIGN KEY (sender) references profile (email)
);
`

const ModelDropQuery = `
DROP TABLE notification;
`
