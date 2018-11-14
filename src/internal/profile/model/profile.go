package model

import "time"

type Profile struct {
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Displayname string    `json:"displayname"`
	Imageurl    string    `json:"imageurl"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE profile (
	username VARCHAR(250) NOT NULL,
	password VARCHAR(250) NOT NULL,
	displayname VARCHAR(250) NOT NULL,
	imageurl Text,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (username)
);
`

const ModelDropQuery = `
DROP TABLE profile;
`
