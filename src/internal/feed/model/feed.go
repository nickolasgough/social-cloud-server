package model

import "time"

type Feed struct {
	Username string    `json:"username"`
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
	username VARCHAR(250) NOT NULL,
	feedname VARCHAR(250) NOT NULL,
	connection VARCHAR(100) NOT NULL,
	joined TIMESTAMP NOT NULL,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username, feedname, datetime, connection, joined),
	FOREIGN KEY (username) references profile (username),
	FOREIGN KEY (connection) references profile (username)
);
`

const ModelDropQuery = `
DROP TABLE feed;
`