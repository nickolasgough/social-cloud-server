package api

import (
	"context"
	"errors"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"time"
	"social-cloud-server/src/internal/util"
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
	Email       string    `json:"email"`
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
	r, ok := request.(*CreateRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a CreateRequest")
	}

	lockIds := []string{"profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	_, err := c.db.ExecStatement(
		c.db.BuildQuery(
			createQuery,
			r.Email,
			r.Password,
			r.DisplayName,
			r.Datetime.Format(time.RFC3339),
		),
	)
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
	email,
	password,
	displayname,
	imageurl,
	defaultfeed,
	datetime
)
VALUES (
	'%s',
	'%s',
	'%s',
	NULL,
	NULL,
	'%s'
);
`
