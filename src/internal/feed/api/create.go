package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/feed/model"
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
	Username string         `json:"username"`
	Feedname string         `json:"feedname"`
	Members  []model.Member `json:"members"`
	Datetime time.Time      `json:"datetime"`
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

	lockIds := []string{"feed"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	for _, member := range r.Members {
		_, err := c.db.ExecStatement(c.db.BuildQuery(createQuery, r.Username, r.Feedname, member.Connection, member.Datetime.Format(time.RFC3339), r.Datetime.Format(time.RFC3339)))
		if err != nil {
			return &CreateResponse{
				Success: false,
			}, err
		}
	}

	return &CreateResponse{
		Success: true,
	}, nil
}

const createQuery = `
INSERT INTO feed (
	username,
	feedname,
	connection,
	joined,
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
