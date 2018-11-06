package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
)

type DismissHandler struct {
	db *database.Database
}

func NewDismissHandler(db *database.Database) *DismissHandler {
	return &DismissHandler{
		db: db,
	}
}

type DismissRequest struct {
	Username string    `json:"username"`
	Sender   string    `json:"sender"`
	Datetime time.Time `json:"datetime"`
}

type DismissResponse struct {
	Success bool `json:"success"`
}

func (c *DismissHandler) Request() endpoint.Request {
	return &DismissRequest{}
}

func (c *DismissHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*DismissRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a DismissRequest")
	}

	_, err := c.db.ExecStatement(c.db.BuildQuery(dismissQuery, r.Username, r.Sender, r.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &DismissResponse{
			Success: false,
		}, err
	}

	return &DismissResponse{
		Success: true,
	}, nil
}

const dismissQuery = `
UPDATE notification
SET dismissed = true
WHERE username = '%s' AND sender = '%s' AND datetime = '%s'
`