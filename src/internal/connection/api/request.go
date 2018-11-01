package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
)

type RequestHandler struct {
	db *database.Database
}

func NewRequestHandler(db *database.Database) *RequestHandler {
	return &RequestHandler{
		db: db,
	}
}

type RequestRequest struct {
	Username   string    `json:"username"`
	Connection string    `json:"connection"`
	Datetime   time.Time `json:"datetime"`
}

type RequestResponse struct {
	Success bool `json:"success"`
}

func (c *RequestHandler) Request() endpoint.Request {
	return &RequestRequest{}
}

func (c *RequestHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	cr, ok := request.(*RequestRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a RequestRequest")
	}

	_, err := c.db.ExecQuery(c.db.BuildQuery(requestQuery, cr.Username, cr.Connection, cr.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &RequestResponse{
			Success: false,
		}, err
	}

	return &RequestResponse{
		Success: true,
	}, nil
}

const requestQuery = `
INSERT INTO notification (
	username,
	type,
	sender,
	dismissed,
	datetime
)
VALUES (
	'%s',
	'connection-request',
	'%s',
	false,
	'%s'
);
`