package model

import "time"

type Profile struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Displayname string    `json:"displayname"`
	Imageurl    string    `json:"imageurl"`
	Datetime    time.Time `json:"datetime"`
}

const ModelCreateQuery = `
CREATE TABLE profile (
	email VARCHAR(250) NOT NULL,
	password VARCHAR(250) NOT NULL,
	displayname VARCHAR(250) NOT NULL,
	imageurl TEXT,
	datetime TIMESTAMP NOT NULL,

	PRIMARY KEY (email)
);
`

const ModelDropQuery = `
DROP TABLE profile;
`
