package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"internal/util"
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
	Imageurl    string `json:"imageurl"`
}

func (c *LoginHandler) Request() endpoint.Request {
	return &LoginRequest{}
}

func (c *LoginHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*LoginRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a LoginRequest")
	}

	lockIds := []string{"profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	result, err := c.db.ExecQuery(c.db.BuildQuery(loginQuery, r.Username, r.Password))
	if err != nil {
		return &LoginResponse{
			Displayname: "",
		}, err
	}

	var lr LoginResponse
	if result.Next() {
		err = result.Scan(&lr.Displayname, &lr.Imageurl)
		if err != nil {
			return &LoginResponse{
				Displayname: "",
				Imageurl: "",
			}, err
		}
	}

	return &lr, nil
}

const loginQuery = `
SELECT
	displayname,
	CASE 
		WHEN imageurl IS NULL THEN ''
		ELSE imageurl
	END
FROM profile
WHERE username = '%s' AND password = '%s';
`
