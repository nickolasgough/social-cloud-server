package model

import (
	"time"
)


type Connection struct {
	Username    string    `json:"username"`
	Connection  string    `json:"connection"`
	Displayname string    `json:"displayname"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE connection (
	username VARCHAR(250) NOT NULL,
	connection VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, datetime),
	FOREIGN KEY (username) references profile (username),
	FOREIGN KEY (connection) references profile (username)
);
`

const ModelDropQuery = `
DROP TABLE connection;
`
