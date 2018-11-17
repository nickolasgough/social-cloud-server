package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/post/model"
)

type ReactHandler struct {
	db *database.Database
}

func NewReactHandler(db *database.Database) *ReactHandler {
	return &ReactHandler{
		db: db,
	}
}

type ReactRequest struct {
	Username  string     `json:"username"`
	Post      model.Post `json:"post"`
	Reaction  string     `json:"reaction"`
	Reacttime time.Time  `json:"reacttime"`
}

type ReactResponse struct {
	Success bool `json:"success"`
}

func (c *ReactHandler) Request() endpoint.Request {
	return &ReactRequest{}
}

func (c *ReactHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*ReactRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ReactRequest")
	}

	var field string
	if r.Reaction == "liked" {
		field = "likes"
	} else {
		field = "dislikes"
	}
	_, err := c.db.ExecStatement(c.db.BuildQuery(reactQuery, field, field, r.Post.Username, r.Post.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &ReactResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(c.db.BuildQuery(notifyQuery, r.Post.Username, r.Reaction, r.Username, r.Reacttime.Format(time.RFC3339)))
	if err != nil {
		return &ReactResponse{
			Success: false,
		}, err
	}

	return &ReactResponse{
		Success: true,
	}, nil
}

const reactQuery = `
UPDATE post
SET %s = %s + 1
WHERE username = '%s' AND datetime = '%s';
`

const notifyQuery = `
INSERT INTO notification (
	username,
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
