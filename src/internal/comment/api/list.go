package api

import (
	"context"
	"errors"
	"time"

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
	Email    string    `json:"email"`
	Datetime time.Time `json:"datetime"`
	Cursor   string    `json:"cursor"`
	Limit    string    `json:"limit"`
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

	results, err := c.db.ExecQuery(
		c.db.BuildQuery(
			listQuery,
			r.Email,
			r.Datetime.Format(time.RFC3339),
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
WHERE co.postemail = '%s' AND co.posttime = '%s'
ORDER BY co.datetime DESC
OFFSET %s
LIMIT %s;
`
