package model

import (
	"time"
)

type Connection struct {
	Email       string    `json:"email"`
	Connection  string    `json:"connection"`
	Displayname string    `json:"displayname"`
	Imageurl    string    `json:"imageurl"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE connection (
	email VARCHAR(250) NOT NULL,
	connection VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (email, datetime),
	FOREIGN KEY (email) references profile (email),
	FOREIGN KEY (connection) references profile (email)
);
`

const ModelDropQuery = `
DROP TABLE connection;
`
