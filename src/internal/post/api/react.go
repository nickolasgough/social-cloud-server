package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/post/model"
	"social-cloud-server/src/internal/util"
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

	lockIds := []string{"post", "reaction", "notification"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	var field string
	if r.Reaction == "liked" {
		field = "likes"
	} else {
		field = "dislikes"
	}

	_, err := c.db.ExecStatement(
		c.db.BuildQuery(
			postQuery,
			field,
			field,
			r.Post.Username,
			r.Post.Datetime.Format(time.RFC3339),
		),
	)
	if err != nil {
		return &ReactResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(
		c.db.BuildQuery(
			reactQuery,
			r.Post.Username,
			r.Post.Datetime.Format(time.RFC3339),
			r.Username,
			r.Reacttime.Format(time.RFC3339),
			r.Reaction,
		),
	)
	if err != nil {
		return &ReactResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(
		c.db.BuildQuery(
			notifyQuery,
			r.Post.Username,
			r.Reaction,
			r.Username,
			r.Reacttime.Format(time.RFC3339),
		),
	)
	if err != nil {
		return &ReactResponse{
			Success: false,
		}, err
	}

	return &ReactResponse{
		Success: true,
	}, nil
}

const postQuery = `
UPDATE post
SET %s = %s + 1
WHERE username = '%s' AND datetime = '%s';
`

const reactQuery = `
INSERT INTO reaction (
	username,
	posttime,
	connection,
	datetime,
	reaction
)
VALUES (
	'%s',
	'%s',
	'%s',
	'%s',
	'%s'
);
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
