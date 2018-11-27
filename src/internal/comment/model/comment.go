package model

import (
	"time"
)

type Comment struct {
	Postemail    string    `json:"postemail"`
	Posttime   time.Time    `json:"posttime"`
	Email string `json:"email"`
	Datetime time.Time `json:"datetime"`
	Avatar Avatar `json:"avatar"`
	Comment     string    `json:"comment"`
}

type Avatar struct {
	Displayname string `json:"displayname"`
	Imageurl    string `json:"imageurl"`
}

const ModelCreateQuery = `
CREATE TABLE comment (
	postemail VARCHAR(250) NOT NULL,
	posttime TIMESTAMP NOT NULL,
	email VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,
	comment TEXT NOT NULL,

	PRIMARY KEY (postemail, posttime, email, datetime),
	FOREIGN KEY (postemail, posttime) REFERENCES post (email, datetime),
	FOREIGN KEY (postemail) REFERENCES profile (email),
	FOREIGN KEY (email) REFERENCES profile (email)
);

CREATE TABLE reaction (
	email VARCHAR(250) NOT NULL,
	posttime TIMESTAMP NOT NULL,
	connection VARCHAR(250) NOT NULL,
	datetime TIMESTAMP NOT NULL,
	reaction VARCHAR(250) NOT NULL,

	PRIMARY KEY (email, posttime, connection, datetime, reaction),
	FOREIGN KEY (email) REFERENCES profile (email),
	FOREIGN KEY (connection) REFERENCES profile (email),
	FOREIGN KEY (email, posttime) REFERENCES post (email, datetime)
);
`

const ModelDropQuery = `
DROP TABLE comment;
`
