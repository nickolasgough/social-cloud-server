package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
)

type CreateHandler struct {
	db *database.Database
}

func NewCreateHandler(db *database.Database) *CreateHandler {
	return &CreateHandler{
		db: db,
	}
}

type CreateRequest struct {
	Username string    `json:"username"`
	Post     string    `json:"post"`
	Datetime time.Time `json:"datetime"`
}

type CreateResponse struct {
	Success bool `json:"success"`
}

func (c *CreateHandler) Request() endpoint.Request {
	return &CreateRequest{}
}

func (c *CreateHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	cr, ok := request.(*CreateRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a CreateRequest")
	}

	_, err := c.db.ExecStatement(c.db.BuildQuery(createQuery, cr.Username, cr.Post, cr.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &CreateResponse{
			Success: false,
		}, err
	}

	return &CreateResponse{
		Success: true,
	}, nil
}

const createQuery = `
INSERT INTO post (
	username,
	post,
	datetime
)
VALUES (
	'%s',
	'%s',
	'%s'
);
`