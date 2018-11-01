package model

import "time"

type Post struct {
	Username string
	Post     string
	Datetime time.Time
}
