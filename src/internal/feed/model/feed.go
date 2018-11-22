package model

import "time"

type Feed struct {
	Email    string    `json:"email"`
	Feedname string    `json:"feedname"`
	Members  []Member  `json:"members"`
	Datetime time.Time `json:"datetime"`
}

type Member struct {
	Connection string    `json:"connection"`
	Datetime   time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE feed (
	email VARCHAR(250) NOT NULL,
	feedname VARCHAR(250) NOT NULL,
	connection VARCHAR(100) NOT NULL,
	joined TIMESTAMP NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (email, feedname, datetime, connection, joined),
	FOREIGN KEY (email) references profile (email),
	FOREIGN KEY (connection) references profile (email)
);
`

const ModelDropQuery = `
DROP TABLE feed;
`
