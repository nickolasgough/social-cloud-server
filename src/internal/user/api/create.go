package api

import (
	"context"
	"errors"

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
	Name string `json:"name"`
}

type CreateResponse struct {
	Name string `json:"name"`
}

func (c *CreateHandler) Request() endpoint.Request {
	return &CreateRequest{}
}

func (c *CreateHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	cr, ok := request.(*CreateRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a CreateRequest")
	}

	err := c.db.ExecQuery(c.db.BuildQuery(createQuery, cr.Name))
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		Name: cr.Name,
	}, nil
}

const createQuery = `
INSERT INTO profile (
	name
)
VALUES (
	'%s'
);
`