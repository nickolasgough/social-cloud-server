package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/connection/model"
	"social-cloud-server/src/internal/util"
)

type ListHandler struct {
	db *database.Database
}

func NewListHandler(db *database.Database) *ListHandler {
	return &ListHandler{
		db: db,
	}
}

type ListRequest struct {
	Username string `json:"username"`
	Cursor   string `json:"cursor"`
	Limit    string `json:"limit"`
}

type ListResponse struct {
	Connections []model.Connection `json:"connections"`
}

func (c *ListHandler) Request() endpoint.Request {
	return &ListRequest{}
}

func (c *ListHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ListRequest")
	}

	lockIds := []string{"connection", "profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	offset := r.Cursor
	if offset == "" {
		offset = "0"
	}
	limit := r.Limit
	if limit == "" {
		limit = "10"
	}
	results, err := c.db.ExecQuery(c.db.BuildQuery(listQuery, r.Username, offset, limit))
	if err != nil {
		return &ListResponse{
			Connections: nil,
		}, err
	}

	var data []model.Connection
	var connection model.Connection
	var datetime string
	for results.Next() {
		err = results.Scan(&connection.Username, &connection.Connection, &connection.Displayname, &datetime)
		if err != nil {
			return &ListResponse{
				Connections: nil,
			}, err
		}

		connection.Datetime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			return ListResponse{
				Connections: nil,
			}, err
		}
		data = append(data, connection)
	}

	return &ListResponse{
		Connections: data,
	}, nil
}

const listQuery = `
SELECT
	c.username,
	c.connection,
	p.displayname,
	c.datetime
FROM connection c
JOIN profile p ON p.username = c.connection
WHERE c.username = '%s'
ORDER BY c.username DESC
OFFSET %s
LIMIT %s;
`