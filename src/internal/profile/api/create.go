package api

import (
	"context"
	"errors"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"time"
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
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	DisplayName string    `json:"displayname"`
	Datetime    time.Time `json:"datetime"`
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

	_, err := c.db.ExecQuery(c.db.BuildQuery(createQuery, cr.Username, cr.Password, cr.DisplayName, cr.Datetime.Format(time.RFC3339)))
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
INSERT INTO profile (
	username,
	password,
	displayname,
	datetime
)
VALUES (
	'%s',
	'%s',
	'%s',
	'%s'
);
`