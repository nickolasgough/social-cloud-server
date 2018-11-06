package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
)

type DeclineHandler struct {
	db *database.Database
}

func NewDeclineHandler(db *database.Database) *DeclineHandler {
	return &DeclineHandler{
		db: db,
	}
}

type DeclineRequest struct {
	Username   string    `json:"username"`
	Connection string    `json:"connection"`
	Datetime   time.Time `json:"datetime"`
}

type DeclineResponse struct {
	Success bool `json:"success"`
}

func (c *DeclineHandler) Request() endpoint.Request {
	return &DeclineRequest{}
}

func (c *DeclineHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*DeclineRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a DeclineRequest")
	}

	_, err := c.db.ExecStatement(c.db.BuildQuery(declineQuery, r.Username, r.Connection))
	if err != nil {
		return &DeclineResponse{
			Success: false,
		}, err
	}

	return &DeclineResponse{
		Success: true,
	}, nil
}

const declineQuery = `
UPDATE notification
SET dismissed = true
WHERE username = '%s' AND sender = '%s' AND type = 'connection-request';
`
