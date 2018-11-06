package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
)

type LoginHandler struct {
	db *database.Database
}

func NewLoginHandler(db *database.Database) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

type LoginRequest struct {
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	DisplayName string    `json:"displayname"`
	Datetime    time.Time `json:"datetime"`
}

type LoginResponse struct {
	Displayname string `json:"displayname"`
}

func (c *LoginHandler) Request() endpoint.Request {
	return &LoginRequest{}
}

func (c *LoginHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*LoginRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a LoginRequest")
	}

	result, err := c.db.ExecQuery(c.db.BuildQuery(loginQuery, r.Username, r.Password))
	if err != nil {
		return &LoginResponse{
			Displayname: "",
		}, err
	}

	var displayname string
	if result.Next() {
		err = result.Scan(&displayname)
		if err != nil {
			return &LoginResponse{
				Displayname: "",
			}, err
		}
	}

	return &LoginResponse{
		Displayname: displayname,
	}, nil
}

const loginQuery = `
SELECT
	displayname
FROM profile
WHERE username = '%s' AND password = '%s';
`