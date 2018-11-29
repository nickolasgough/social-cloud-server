package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
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
	Postemail string    `json:"postemail"`
	Posttime  time.Time `json:"posttime"`
	Email     string    `json:"email"`
	Datetime  time.Time `json:"datetime"`
	Comment   string    `json:"comment"`
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

	lockIds := []string{"comment"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	_, err := c.db.ExecStatement(
		c.db.BuildQuery(
			createQuery,
			r.Postemail,
			r.Posttime.Format(time.RFC3339),
			r.Email,
			r.Comment,
			r.Datetime.Format(time.RFC3339),
		),
	)
	if err != nil {
		return &CreateResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(c.db.BuildQuery(postQuery, r.Postemail, r.Posttime.Format(time.RFC3339)))
	if err != nil {
		return &CreateResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(
		c.db.BuildQuery(
			notifyQuery,
			r.Postemail,
			"comment",
			r.Email,
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
INSERT INTO comment (
	postemail,
	posttime,
	email,
	comment,
	datetime
)
VALUES (
	'%s',
	'%s',
	'%s',
	'%s',
	'%s'
);
`

const postQuery = `
UPDATE post
SET comments = comments + 1
WHERE email = '%s' AND datetime = '%s';
`

const notifyQuery = `
INSERT INTO notification (
	email,
	type,
	sender,
	dismissed,
	datetime
)
VALUES (
	'%s',
	'post-%s',
	'%s',
	false,
	'%s'
)
`
