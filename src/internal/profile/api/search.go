package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
	"social-cloud-server/src/internal/profile/model"
)

type SearchHandler struct {
	db *database.Database
}

func NewSearchHandler(db *database.Database) *SearchHandler {
	return &SearchHandler{
		db: db,
	}
}

type SearchRequest struct {
	Email string `json:"email"`
	Query string `json:"query"`
}

type SearchResponse struct {
	Users []model.User `json:"users"`
}

func (c *SearchHandler) Request() endpoint.Request {
	return &SearchRequest{}
}

func (c *SearchHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*SearchRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a SearchRequest")
	}

	lockIds := []string{"profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	result, err := c.db.ExecQuery(c.db.BuildQuery(searchQuery, r.Email, r.Email, r.Email, r.Query))
	if err != nil {
		return &SearchResponse{
			Users: nil,
		}, err
	}

	var users []model.User
	var user model.User
	var datetime string
	var ccount int
	var rcount int
	for result.Next() {
		err = result.Scan(&user.Email, &user.Displayname, &user.Imageurl, &ccount, &rcount, &datetime)
		if err != nil {
			return &SearchResponse{
				Users: nil,
			}, err
		}

		user.Datetime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			return &SearchResponse{
				Users: nil,
			}, err
		}

		user.Connected = ccount > 0 || rcount > 0
		users = append(users, user)
	}

	return &SearchResponse{
		Users: users,
	}, nil
}

const searchQuery = `
SELECT
	p.email,
	p.displayname,
	CASE 
		WHEN p.imageurl IS NULL THEN ''
		ELSE p.imageurl
	END,
	(SELECT
		COUNT(c.connection)
	FROM connection c
	WHERE c.email = '%s' AND c.connection = p.email),
	(SELECT
		COUNT(n.email)
	FROM notification n
	WHERE n.sender = '%s' AND n.email = p.email AND n.type = 'connection-request' AND n.dismissed = false),
	p.datetime
FROM profile p
WHERE p.email != '%s' AND p.displayname LIKE '%%%s%%';
`
