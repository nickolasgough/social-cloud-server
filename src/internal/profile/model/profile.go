package model


type Profile struct {
	Username    string
	Password    string
	Displayname string
}

const ModelCreateQuery = `
CREATE TABLE profile (
	username VARCHAR(250) NOT NULL,
	password VARCHAR(250) NOT NULL,
	displayname VARCHAR(250) NOT NULL,

	PRIMARY KEY (username)
);
`

const ModelDropQuery = `
DROP TABLE profile;
`
