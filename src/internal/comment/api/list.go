package api

import (
	"context"
	"errors"
	"time"
	"fmt"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/comment/model"
	"social-cloud-server/src/internal/util"
)

type ListHandler struct {
	db *database.Database
}

func NewListHandler(db *database.Database) *ListHandler {
	return &ListHandler{
		db: db,
	}
}

type ListRequest struct {
	Postemail string    `json:"postemail"`
	Posttime  time.Time `json:"posttime"`
	Email     string    `json:"email"`
	Feedname  string    `json:"feedname"`
	Cursor    string    `json:"cursor"`
	Limit     string    `json:"limit"`
}

type ListResponse struct {
	Comments []model.Comment `json:"comments"`
}

func (c *ListHandler) Request() endpoint.Request {
	return &ListRequest{}
}

func (c *ListHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ListRequest")
	}

	lockIds := []string{"comment", "profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	offset := r.Cursor
	if offset == "" {
		offset = "0"
	}
	limit := r.Limit
	if limit == "" {
		limit = "25"
	}

	var conditionQuery string
	if r.Feedname != "" {
		conditionQuery = fmt.Sprintf(
			feedCondition,
			r.Postemail,
			r.Posttime.Format(time.RFC3339),
			r.Email,
			r.Email,
			r.Feedname,
		)
	} else {
		conditionQuery = fmt.Sprintf(
			userCondition,
			r.Postemail,
			r.Posttime.Format(time.RFC3339),
		)
	}

	results, err := c.db.ExecQuery(
		c.db.BuildQuery(
			listQuery,
			conditionQuery,
			offset,
			limit,
		),
	)
	if err != nil {
		return &ListResponse{
			Comments: nil,
		}, err
	}

	var comments []model.Comment
	var comment model.Comment
	var avator model.Avatar
	var datetime string
	for results.Next() {
		err = results.Scan(
			&comment.Email,
			&comment.Comment,
			&avator.Displayname,
			&avator.Imageurl,
			&datetime,
		)
		if err != nil {
			return &ListResponse{
				Comments: nil,
			}, err
		}
		comment.Avatar = avator

		comment.Datetime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			return ListResponse{
				Comments: nil,
			}, err
		}

		comments = append(comments, comment)
	}

	return &ListResponse{
		Comments: comments,
	}, nil
}

const listQuery = `
SELECT
	co.email,
	co.comment,
	pr.displayname,
	pr.imageurl,
	co.datetime
FROM comment co
JOIN profile pr ON pr.email = co.email
WHERE %s
ORDER BY co.datetime DESC
OFFSET %s
LIMIT %s;
`

const feedCondition = `
co.postemail = '%s' AND co.posttime = '%s' AND (co.email = '%s' OR co.email IN (
	SELECT
		fd.connection
	FROM feed fd
	WHERE fd.email = '%s' AND fd.feedname = '%s'))
`

const userCondition = `
co.postemail = '%s' AND co.posttime = '%s'
`
